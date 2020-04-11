package main

//go:generate rice embed-go

import (
	"context"
	"database/sql"

	rice "github.com/GeertJohan/go.rice"
	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/andrewstucki/web-app-tools/go/server"
	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"

	"example/models"
	"example/roles"
	"example/routes"
)

// This gets a user based off of a JWT identity token
func userFromClaims(ctx context.Context, claims *verifier.StandardClaims) (*models.User, error) {
	user, err := models.Users(models.UserWhere.GoogleID.EQ(claims.Subject)).One(ctx, server.Tx(ctx))
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// This gets a user based off of an API token
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

var config = server.Config{
	// This contains the path to the migrations
	Migrations: rice.MustFindBox("./migrations"),
	// This contains the path to the frontend assets
	Assets: rice.MustFindBox("./frontend/build"),
	// This method is invoked after the database is initialized and migrated and various
	// internal auth handlers are set up, it's where the main route initialization code
	// should be
	Setup: func(config *server.SetupConfig) {
		roles.Register()
		config.Router.Route("/v1", routes.NewV1Handler(config.Logger, config.Render).Register)
	},
	// This callback is used to actually return the currently logged in user and allows
	// us to invoked the server.CurrentUser(ctx) method to get the user
	GetCurrentUser: func(ctx context.Context, claimsOrToken *server.ClaimsOrToken) (interface{}, error) {
		if claimsOrToken.Claims != nil {
			return userFromClaims(ctx, claimsOrToken.Claims)
		}
		return userFromToken(ctx, claimsOrToken.Token)
	},
	// This gets invoked the first time anyone ever logs into the system
	// it's useful for setting up an admin user
	OnFirstUser: func(ctx context.Context, claims *verifier.StandardClaims) error {
		user := models.User{Email: claims.Email, GoogleID: claims.Subject}
		if err := user.Upsert(ctx, server.Tx(ctx), false, nil, boil.Infer(), boil.Infer()); err != nil {
			return err
		}
		return security.SetRole(ctx, roles.SuperAdminRole, user.ID)
	},
	// This gets called every subsequent log in, it's useful for inserting
	// a user if they don't already exist
	OnLogin: func(ctx context.Context, claims *verifier.StandardClaims) error {
		user := models.User{Email: claims.Email, GoogleID: claims.Subject}
		return user.Upsert(ctx, server.Tx(ctx), false, nil, boil.Infer(), boil.Infer())
	},
}

func main() {
	server.RunServer(config)
}
