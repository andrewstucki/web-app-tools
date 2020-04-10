package main

//go:generate rice embed-go

import (
	"context"
	"database/sql"

	"example/models"
	"example/roles"
	"example/routes"

	rice "github.com/GeertJohan/go.rice"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/andrewstucki/web-app-tools/go/server"
)

func userFromClaims(ctx context.Context, claims *verifier.StandardClaims) (*models.User, error) {
	user, err := models.Users(models.UserWhere.GoogleID.EQ(claims.Subject)).One(ctx, server.Tx(ctx))
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func userFromToken(ctx context.Context, token string) (*models.User, error) {
	if apiToken, err := uuid.FromString(token); err == nil {
		user, err := models.Users(qm.InnerJoin("user_tokens on user_tokens.user_id = users.id"), models.UserTokenWhere.ID.EQ(apiToken)).One(ctx, server.Tx(ctx))
		if err != nil && err == sql.ErrNoRows {
			return nil, nil
		}
		return user, err
	}
	// we don't have a valid uuid
	return nil, nil
}

func main() {
	server.RunServer(server.Config{
		Migrations: rice.MustFindBox("./migrations"),
		Assets:     rice.MustFindBox("./frontend/build"),
		Setup: func(config *server.SetupConfig) {
			roles.Register()
			config.Router.Route("/v1", routes.NewV1Handler(config.Logger, config.Render).Register)
		},
		GetCurrentUser: func(ctx context.Context, claimsOrToken *server.ClaimsOrToken) (interface{}, error) {
			if claimsOrToken.Claims != nil {
				return userFromClaims(ctx, claimsOrToken.Claims)
			}
			return userFromToken(ctx, claimsOrToken.Token)
		},
		OnFirstUser: func(ctx context.Context, claims *verifier.StandardClaims) error {
			user := models.User{Email: claims.Email, GoogleID: claims.Subject}
			if err := user.Upsert(ctx, server.Tx(ctx), false, nil, boil.Infer(), boil.Infer()); err != nil {
				return err
			}
			return security.SetRole(ctx, roles.SuperAdminRole, user.ID)
		},
		OnLogin: func(ctx context.Context, claims *verifier.StandardClaims) error {
			user := models.User{Email: claims.Email, GoogleID: claims.Subject}
			return user.Upsert(ctx, server.Tx(ctx), false, nil, boil.Infer(), boil.Infer())
		},
	})
}
