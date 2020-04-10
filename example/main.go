package main

//go:generate rice embed-go
//go:generate sqlboiler --wipe -c .sqlboiler.toml psql

import (
	"context"
	"database/sql"
	"example/models"

	rice "github.com/GeertJohan/go.rice"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	"github.com/andrewstucki/web-app-tools/go/server"
	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"
)

var (
	logger zerolog.Logger
	render common.Renderer

	persistUserQuery = `
	INSERT INTO users (email, google_id)
		VALUES ($1, $2)
	ON CONFLICT DO NOTHING;
	`
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
			_, err := config.DB.Exec(persistUserQuery, claims.Email, claims.Subject)
			return err
		},
	})
}
