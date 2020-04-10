package main

//go:generate rice embed-go
//go:generate sqlboiler --no-hooks --no-rows-affected --no-tests --wipe -c .sqlboiler.toml psql

import (
	"context"
	"database/sql"
	"example/models"

	rice "github.com/GeertJohan/go.rice"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	"github.com/andrewstucki/web-app-tools/go/server"
	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"
)

func main() {
	server.RunServer(server.Config{
		Migrations: rice.MustFindBox("./migrations"),
		Assets:     rice.MustFindBox("./frontend/build"),
		HostPort:   ":3456",
		Domains:    []string{"gpmail.org"},
		Setup: func(config *server.SetupConfig) {
			boil.SetDB(config.DB)
			logger = config.Logger
			render = config.Render
			config.Router.Get("/v1/me", me)
		},
		GetCurrentUser: func(ctx context.Context, queryer sqlContext.QueryContext, claimsOrToken *server.ClaimsOrToken) (interface{}, error) {
			var user *models.User
			var err error
			if claimsOrToken.Claims != nil {
				claims := claimsOrToken.Claims
				user, err = models.Users(models.UserWhere.GoogleID.EQ(claims.Subject)).One(ctx, queryer)
			} else {
				token := claimsOrToken.Token
				user, err = models.Users(qm.Load(models.UserRels.UserTokens, models.UserTokenWhere.ID.EQ(token))).One(ctx, queryer)
			}
			if err != nil && err == sql.ErrNoRows {
				return nil, nil
			}
			return user, err
		},
		OnLogin: func(config *server.SetupConfig, claims *verifier.StandardClaims) error {
			user := models.User{Email: claims.Email, GoogleID: claims.Subject}
			return user.Upsert(context.Background(), config.DB.DB, false, nil, boil.Infer(), boil.Infer())
		},
	})
}
