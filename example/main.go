package main

// go:generate rice embed-go
// go:generate sqlboiler --no-hooks --no-rows-affected --no-tests --wipe -c .sqlboiler.toml psql

import (
	"context"
	"database/sql"
	"example/models"

	rice "github.com/GeertJohan/go.rice"
	"github.com/volatiletech/sqlboiler/boil"

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
		GetCurrentUser: func(ctx context.Context, queryer sqlContext.QueryContext, claims *verifier.StandardClaims) (interface{}, error) {
			user, err := models.Users(models.UserWhere.GoogleID.EQ(claims.Subject)).One(ctx, queryer)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, nil
				}
				return nil, err
			}
			return user, nil
		},
		OnLogin: func(config *server.SetupConfig, claims *verifier.StandardClaims) error {
			user := models.User{Email: claims.Email, GoogleID: claims.Subject}
			return user.Upsert(context.Background(), config.DB.DB, false, nil, boil.Infer(), boil.Infer())
		},
	})
}
