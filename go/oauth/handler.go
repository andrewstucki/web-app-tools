package server

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"github.com/andrewstucki/web-app-tools/go/oauth/callbacks"
	"github.com/andrewstucki/web-app-tools/go/oauth/state"
	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
)

const successTemplate = `
<!doctype html>
  <script>try { localStorage.setItem("__google_id", "{{.}}") } finally { window.location = "/" }</script>
</html>
`

var (
	defaultSuccessTemplate *template.Template
	contextKey             = struct{}{}
)

func init() {
	defaultSuccessTemplate = template.Must(template.New("__oauth__success").Parse(successTemplate))
}

// Handler handles oauth2 authentication requests.
type Handler struct {
	*http.ServeMux

	config           *oauth2.Config
	url              string
	timeout          time.Duration
	tokenManager     TokenManager
	verifier         *verifier.Verifier
	callbacks        Callbacks
	secretKey        string
	allowedRedirects []string
	logger           zerolog.Logger
}

// New creates a new handler based on the given config.
func New(config *Config) (*Handler, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	timeout := config.ClientTimeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	tokenVerifier := config.Verifier
	if tokenVerifier == nil {
		tokenVerifier = verifier.NewVerifier()
	}

	tokenManager := config.TokenManager
	if tokenManager == nil {
		tokenManager = state.NewMemoryTokenManager()
	}

	tokenCallbacks := config.Callbacks
	if tokenCallbacks == nil {
		tokenCallbacks = callbacks.NewLocalStorageCallbacks()
	}

	allowedRedirects := config.AllowedRedirects
	if len(allowedRedirects) == 0 {
		allowedRedirects = []string{"/"}
	}

	logger := zerolog.Nop()
	if config.Logger != nil {
		logger = *config.Logger
	}

	h := &Handler{
		ServeMux: http.NewServeMux(),
		config: &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.MountURL + "/callback",
			Endpoint:     google.Endpoint,
			Scopes:       []string{"openid", "profile", "email"},
		},
		url:              config.MountURL,
		callbacks:        tokenCallbacks,
		timeout:          timeout,
		tokenManager:     tokenManager,
		verifier:         tokenVerifier,
		secretKey:        config.SecretKey,
		allowedRedirects: allowedRedirects,
		logger:           logger,
	}

	h.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		startPath := config.mountURL.Path
		callbackPath := config.mountURL.Path + "/callback"

		if req.Method != "GET" {
			http.NotFound(w, req)
			return
		}
		if req.URL.Path == startPath {
			h.handleBegin(w, req)
			return
		}
		if req.URL.Path == callbackPath {
			h.handleEnd(w, req)
			return
		}
		http.NotFound(w, req)
		return
	})

	return h, nil
}

func (h *Handler) handleBegin(w http.ResponseWriter, r *http.Request) {
	disableCaching(w)

	location := r.URL.Query().Get("redirect")
	if location == "" {
		location = "/"
	}

	found := false
	for _, whitelisted := range h.allowedRedirects {
		if whitelisted == location {
			found = true
			break
		}
	}

	if !found {
		h.logger.Warn().Err(ErrInvalidRedirect).Msgf("attempted redirect to '%s'", location)
		h.callbacks.OnError(w, ErrInvalidRedirect)
		return
	}

	state, err := h.generateState(location)
	if err != nil {
		h.logger.Warn().Err(err).Msg(MessageStateGenerationFailed)
		h.callbacks.OnError(w, errors.Wrap(err, MessageStateGenerationFailed))
		return
	}

	http.Redirect(w, r, h.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), http.StatusFound)
}

func (h *Handler) handleEnd(w http.ResponseWriter, r *http.Request) {
	disableCaching(w)

	queryState := r.URL.Query().Get("state")
	if queryState == "" {
		h.logger.Warn().Err(ErrInvalidStateValue).Msg("state is empty")
		h.callbacks.OnError(w, ErrInvalidStateValue)
		return
	}

	location, err := h.validateState(queryState)
	if err != nil {
		h.logger.Warn().Err(err).Msg("state failed to validate")
		if err == ErrInvalidStateValue {
			h.callbacks.OnError(w, ErrInvalidStateValue)
		} else {
			h.callbacks.OnError(w, errors.Wrap(err, ErrInvalidStateValue.Error()))
		}
		return
	}

	queryCode := r.URL.Query().Get("code")
	if queryCode == "" {
		h.logger.Warn().Err(ErrInvalidCodeValue).Msg("code is empty")
		h.callbacks.OnError(w, ErrInvalidCodeValue)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	token, err := h.config.Exchange(ctx, queryCode)
	if err != nil {
		h.logger.Warn().Err(err).Msg(MessageExchangeFailed)
		h.callbacks.OnError(w, errors.Wrap(err, MessageExchangeFailed))
		return
	}

	rawToken, err := h.getClaimsAndCacheToken(r.Context(), token)
	if err != nil {
		h.callbacks.OnInvalidToken(w, err)
		return
	}

	h.callbacks.OnSuccess(w, location, rawToken)
}

// AuthenticationMiddleware provides a mechanism for validating tokens passed
// in Authorization headers
func (h *Handler) AuthenticationMiddleware(unauthorizedHandler func(w http.ResponseWriter)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(auth) != 2 || strings.ToLower(auth[0]) != "bearer" {
				h.logger.Info().Msg("no authorization headers present")
				unauthorizedHandler(w)
				return
			}

			var rawToken string
			tokenClaims := &verifier.StandardClaims{}
			// bad claims == bad token
			if err := h.verifier.VerifyIDToken(auth[1], tokenClaims); err != nil {
				h.logger.Warn().Err(err).Msg("failed to verify token")
				unauthorizedHandler(w)
				return
			}

			expiration := time.Unix(tokenClaims.ExpiresAt, 0)
			if time.Until(expiration) < 10*time.Minute {
				// refresh the token, if anything apart from our hook
				// fails, then just don't do anything until the next request
				serialized, err := h.tokenManager.Get(r.Context(), tokenClaims.Subject)
				if err != nil {
					h.logger.Warn().Err(err).Msg("failed to retrieve token from manager")
					goto SET_CONTEXT
				}

				token := &oauth2.Token{}
				// just bail if we return an error
				if err := json.Unmarshal([]byte(serialized), token); err != nil {
					h.logger.Warn().Err(err).Msg("failed to unmarshal token")
					goto SET_CONTEXT
				}

				// this is a nasty hack to force token refreshing
				token.Expiry = (time.Time{}).Add(1 * time.Second)

				source := h.config.TokenSource(oauth2.NoContext, token)
				refreshed, err := source.Token()
				if err != nil {
					h.logger.Warn().Err(err).Msg("failed to refresh token")
					goto SET_CONTEXT
				}

				rawToken, err = h.getClaimsAndCacheToken(r.Context(), refreshed)
				if err != nil {
					goto SET_CONTEXT
				}

				if err := h.callbacks.OnRefresh(w, rawToken); err != nil {
					h.logger.Warn().Err(err).Msg("refresh handler failed")
				}

			}

		SET_CONTEXT:
			ctx := context.WithValue(r.Context(), contextKey, tokenClaims)
			next.ServeHTTP(w, r.Clone(ctx))
		})
	}
}

// Claims returns claims if they exist on the context
func (h *Handler) Claims(ctx context.Context) *verifier.StandardClaims {
	claims := ctx.Value(contextKey)
	if claims == nil {
		return nil
	}
	return claims.(*verifier.StandardClaims)
}

// MustClaims panics if no claims exist on the context
func (h *Handler) MustClaims(ctx context.Context) *verifier.StandardClaims {
	claims := ctx.Value(contextKey)
	if claims == nil {
		panic("claims not found on context")
	}
	return claims.(*verifier.StandardClaims)
}

// any errors here are going to result in an ErrInvalidToken above
func (h *Handler) getClaimsAndCacheToken(ctx context.Context, token *oauth2.Token) (string, error) {
	if !token.Valid() {
		h.logger.Warn().Err(ErrInvalidToken).Msg("token failed validation")
		return "", ErrInvalidToken
	}

	// convert the id_token parameter
	extraToken := token.Extra("id_token")
	if extraToken == nil {
		h.logger.Warn().Err(ErrInvalidToken).Msg("no id_token value found")
		return "", ErrInvalidToken
	}
	idToken, ok := extraToken.(string)
	if !ok {
		h.logger.Warn().Err(ErrInvalidToken).Msg("id_token of the wrong type")
		return "", ErrInvalidToken
	}
	tokenClaims := &verifier.StandardClaims{}
	if err := h.verifier.VerifyIDToken(idToken, tokenClaims); err != nil {
		h.logger.Warn().Err(err).Msg("token verification failed")
		return "", err
	}

	// use the token manager to store the token serialized as JSON
	serialized, err := json.Marshal(token)
	if err != nil {
		h.logger.Warn().Err(err).Msg("json marshaling failed")
		return "", err
	}
	if err := h.tokenManager.Set(ctx, tokenClaims.Subject, string(serialized)); err != nil {
		h.logger.Warn().Err(err).Msg("failed to write token to manager")
		return "", err
	}
	return idToken, nil
}

type stateClaims struct {
	jwt.StandardClaims
	Location string `json:"location"`
}

func (h *Handler) generateState(location string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stateClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
		Location: location,
	})
	return token.SignedString([]byte(h.secretKey))
}

func (h *Handler) validateState(token string) (string, error) {
	claims := &stateClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			h.logger.Warn().Msg("wrong signing method for token")
			return nil, ErrInvalidStateValue
		}
		return []byte(h.secretKey), nil
	})
	if err != nil {
		return "", err
	}
	return claims.Location, nil
}
