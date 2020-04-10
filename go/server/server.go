package server

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/go-chi/chi"
	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/oauth"
	"github.com/andrewstucki/web-app-tools/go/oauth/callbacks"
	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/andrewstucki/web-app-tools/go/server/middleware"
	"github.com/andrewstucki/web-app-tools/go/sql"
	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"
	sqlMiddleware "github.com/andrewstucki/web-app-tools/go/sql/middleware"
	"github.com/andrewstucki/web-app-tools/go/sql/migrator"
	sqlSecurity "github.com/andrewstucki/web-app-tools/go/sql/security"
	"github.com/andrewstucki/web-app-tools/go/sql/state"
)

func init() {
	// ignore the error if no .env file is found
	godotenv.Load()
}

type assetServer struct {
	*rice.Box
}

func (a *assetServer) Open(name string) (http.File, error) {
	if file, err := a.Box.Open(name); err == nil {
		return file, err
	}
	return a.Box.Open("index.html")
}

func newAssetServer(box *rice.Box) http.FileSystem {
	return &assetServer{
		Box: box,
	}
}

// SetupConfig provides the configuration for the setup phase
type SetupConfig struct {
	Render  common.Renderer
	Router  chi.Router
	DB      *sqlx.DB
	Logger  zerolog.Logger
	Handler *oauth.Handler
}

// Config provides the configuration for the server
type Config struct {
	Migrations     *rice.Box
	Assets         *rice.Box
	HostPort       string
	DatabaseURL    string
	ClientID       string
	ClientSecret   string
	BaseURL        string
	SecretKey      string
	Domains        []string
	Setup          func(config *SetupConfig)
	GetCurrentUser func(ctx context.Context, claimsOrToken *ClaimsOrToken) (interface{}, error)
	OnFirstUser    func(ctx context.Context, claims *verifier.StandardClaims) error
	OnLogin        func(ctx context.Context, claims *verifier.StandardClaims) error
}

type wrappedCallbacks struct {
	*callbacks.LocalStorageCallbacks
	config      *SetupConfig
	initialHook func(ctx context.Context, claims *verifier.StandardClaims) error
	hook        func(ctx context.Context, claims *verifier.StandardClaims) error

	mutex       sync.Mutex
	initialized bool
}

func (c *wrappedCallbacks) checkAndInitialize(claims *verifier.StandardClaims) error {
	tx, ctx, err := sqlContext.StartTx(context.Background(), c.config.DB)
	if err != nil {
		return err
	}
	if err := tx.GetContext(ctx, &c.initialized, "SELECT initialized FROM site_settings"); err != nil {
		c.config.Logger.Error().Err(err).Msg("error checking site settings")
		tx.Rollback()
		return err
	}
	if !c.initialized {
		if err := c.initialHook(ctx, claims); err != nil {
			c.config.Logger.Error().Err(err).Msg("error running first user callback")
			tx.Rollback()
			return err
		}
		if _, err := tx.ExecContext(ctx, "UPDATE site_settings SET initialized = TRUE"); err != nil {
			c.config.Logger.Error().Err(err).Msg("error updating site settings")
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		c.config.Logger.Error().Err(err).Msg("error committing first user transaction")
		tx.Rollback()
		return err
	}
	return nil
}

func (c *wrappedCallbacks) callHook(claims *verifier.StandardClaims) error {
	tx, ctx, err := sqlContext.StartTx(context.Background(), c.config.DB)
	if err != nil {
		return err
	}
	if err := c.hook(ctx, claims); err != nil {
		c.config.Logger.Error().Err(err).Msg("error running login callback")
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		c.config.Logger.Error().Err(err).Msg("error committing login transaction")
		tx.Rollback()
		return err
	}
	return nil
}

func (c *wrappedCallbacks) OnSuccess(w http.ResponseWriter, location, raw string, claims *verifier.StandardClaims) {
	ran := false
	if c.initialHook != nil {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		if !c.initialized {
			if err := c.checkAndInitialize(claims); err != nil {
				c.LocalStorageCallbacks.OnError(w, err)
				return
			}
			c.initialized = true
			ran = true
		}
	}
	if !ran && c.hook != nil {
		if err := c.callHook(claims); err != nil {
			c.LocalStorageCallbacks.OnError(w, err)
			return
		}
	}
	c.LocalStorageCallbacks.OnSuccess(w, location, raw, claims)
}

// RunServer runs a server with the specified config
func RunServer(config Config) {
	logger := zerolog.New(os.Stdout)

	if config.Migrations == nil {
		logger.Fatal().Msg("must specify migrations")
	}

	dbURL := config.DatabaseURL
	if dbURL == "" {
		dbURL = os.Getenv("POSTGRES_URL")
	}

	migrator, err := migrator.NewBoxMigrator(config.Migrations, dbURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to initialize migrator")
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal().Err(err).Msg("failed to run migrations")
	}

	db, err := sql.Connect(dbURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}

	security.RegisterManager(sqlSecurity.NewNamespaceManager(db))

	render := common.NewJSONRenderer()
	router := chi.NewRouter()
	setupConfig := &SetupConfig{
		Render: render,
		DB:     db,
		Logger: logger,
	}

	handler, err := initializeOAuth(setupConfig, config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize oauth handler")
	}

	router.Mount("/oauth", handler)
	router.Route("/api", func(router chi.Router) {
		router.Use(
			middleware.RequestLogger(logger),
			middleware.Recoverer(render, logger),
			handler.AuthenticationMiddleware(false, func(w http.ResponseWriter) {
				render.Error(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}),
			tokenUser(handler),
			sqlMiddleware.Transaction(db, render, logger),
			currentUser(handler, render, logger, config.GetCurrentUser),
		)
		setupConfig.Router = router
		setupConfig.Handler = handler
		config.Setup(setupConfig)
	})
	if config.Assets != nil {
		router.Handle("/*", http.FileServer(newAssetServer(config.Assets)))
	}

	hostPort := config.HostPort
	if hostPort == "" {
		hostPort = os.Getenv("HOST_PORT")
	}

	if err := http.ListenAndServe(hostPort, router); err != nil {
		logger.Fatal().Err(err).Msg("failed to run server")
	}
}

func initializeOAuth(setup *SetupConfig, config Config) (*oauth.Handler, error) {
	verifier := verifier.NewVerifier()
	domains := config.Domains
	if len(config.Domains) == 0 {
		if envDomains := strings.TrimSpace(os.Getenv("GOOGLE_DOMAINS")); envDomains != "" {
			domains = strings.Split(envDomains, ",")
			for i, domain := range domains {
				domains[i] = strings.TrimSpace(domain)
			}
		}
	}
	if len(domains) > 0 {
		verifier = verifier.WithDomains(domains...)
	}

	clientID := config.ClientID
	if clientID == "" {
		clientID = os.Getenv("GOOGLE_CLIENT_ID")
	}
	clientSecret := config.ClientSecret
	if clientSecret == "" {
		clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	}
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("BASE_URL")
	}
	secretKey := config.SecretKey
	if secretKey == "" {
		secretKey = os.Getenv("JWT_SECRET")
	}

	return oauth.New(&oauth.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		MountURL:     baseURL + "/oauth",
		SecretKey:    secretKey,
		Verifier:     verifier,
		TokenManager: state.NewTokenManager(setup.DB),
		Callbacks: &wrappedCallbacks{
			LocalStorageCallbacks: callbacks.NewLocalStorageCallbacks(),
			config:                setup,
			initialHook:           config.OnFirstUser,
			hook:                  config.OnLogin,
		},
		Logger: &setup.Logger,
	})
}
