package main

import (
	"context"
	"testing"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	"github.com/andrewstucki/web-app-tools/go/security"
	"github.com/andrewstucki/web-app-tools/go/security/memory"
	"github.com/andrewstucki/web-app-tools/go/server"
	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/boil"

	"example/models"
	"example/roles"
)

const postgresURL = "postgres://postgres:postgres@localhost:5434/example-test?sslmode=disable"

func randomUser() *models.User {
	return &models.User{
		ID:        uuid.NewV4(),
		Email:     randomdata.Email(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		GoogleID:  randomdata.StringNumber(16, ""),
	}
}

func testInTransaction(t *testing.T, inner func(ctx context.Context)) {
	db := sqlx.MustConnect("postgres", postgresURL)
	defer db.Close()

	tx, ctx, err := sqlContext.StartTx(context.Background(), db)
	require.NoError(t, err)
	inner(ctx)
	require.NoError(t, tx.Rollback())
}

func TestGetCurrentUser_InvalidToken(t *testing.T) {
	testInTransaction(t, func(ctx context.Context) {
		returned, err := config.GetCurrentUser(ctx, &server.ClaimsOrToken{
			Token: "1",
		})
		require.NoError(t, err)
		require.Nil(t, returned)
	})
}

func TestGetCurrentUser_ValidTokenFound(t *testing.T) {
	testInTransaction(t, func(ctx context.Context) {
		user := randomUser()
		token := &models.UserToken{}
		require.NoError(t, user.Insert(ctx, sqlContext.FromContext(ctx), boil.Infer()))
		require.NoError(t, user.AddUserTokens(ctx, sqlContext.FromContext(ctx), true, token))
		returned, err := config.GetCurrentUser(ctx, &server.ClaimsOrToken{
			Token: token.ID.String(),
		})
		require.NoError(t, err)
		require.NotNil(t, returned)
		require.Equal(t, user.ID, returned.(*models.User).ID)
	})
}

func TestGetCurrentUser_ValidTokenNotFound(t *testing.T) {
	testInTransaction(t, func(ctx context.Context) {
		user := randomUser()
		token := &models.UserToken{}
		require.NoError(t, user.Insert(ctx, sqlContext.FromContext(ctx), boil.Infer()))
		require.NoError(t, user.AddUserTokens(ctx, sqlContext.FromContext(ctx), true, token))
		returned, err := config.GetCurrentUser(ctx, &server.ClaimsOrToken{
			Token: uuid.NewV4().String(),
		})
		require.NoError(t, err)
		require.Nil(t, returned)
	})
}

func TestGetCurrentUser_ClaimsNotFound(t *testing.T) {
	testInTransaction(t, func(ctx context.Context) {
		user := randomUser()
		require.NoError(t, user.Insert(ctx, sqlContext.FromContext(ctx), boil.Infer()))
		returned, err := config.GetCurrentUser(ctx, &server.ClaimsOrToken{
			Claims: &verifier.StandardClaims{
				Subject: randomdata.StringNumber(16, ""),
			},
		})
		require.NoError(t, err)
		require.Nil(t, returned)
	})
}

func TestGetCurrentUser_ClaimsFound(t *testing.T) {
	testInTransaction(t, func(ctx context.Context) {
		user := randomUser()
		require.NoError(t, user.Insert(ctx, sqlContext.FromContext(ctx), boil.Infer()))
		returned, err := config.GetCurrentUser(ctx, &server.ClaimsOrToken{
			Claims: &verifier.StandardClaims{
				Subject: user.GoogleID,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, returned)
		require.Equal(t, user.ID, returned.(*models.User).ID)
	})
}

func TestOnLogin(t *testing.T) {
	testInTransaction(t, func(ctx context.Context) {
		security.RegisterManager(memory.NewNamespaceManager())
		security.Register(roles.SuperAdminRole)

		claims := &verifier.StandardClaims{
			Email:   randomdata.Email(),
			Subject: randomdata.StringNumber(16, ""),
		}
		err := config.OnLogin(ctx, claims)
		require.NoError(t, err)

		user, err := models.Users(models.UserWhere.GoogleID.EQ(claims.Subject)).One(ctx, sqlContext.FromContext(ctx))
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, claims.Email, user.Email)

		policies, err := security.WithUser(user.ID).Policies(ctx)
		require.NoError(t, err)
		require.Empty(t, policies)
	})
}

func TestOnFirstUser(t *testing.T) {
	testInTransaction(t, func(ctx context.Context) {
		security.RegisterManager(memory.NewNamespaceManager())
		security.Register(roles.SuperAdminRole)

		claims := &verifier.StandardClaims{
			Email:   randomdata.Email(),
			Subject: randomdata.StringNumber(16, ""),
		}
		err := config.OnFirstUser(ctx, claims)
		require.NoError(t, err)

		user, err := models.Users(models.UserWhere.GoogleID.EQ(claims.Subject)).One(ctx, sqlContext.FromContext(ctx))
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, claims.Email, user.Email)

		policies, err := security.WithUser(user.ID).Policies(ctx)
		require.NoError(t, err)
		require.ElementsMatch(t, roles.SuperAdminRole.Policies, policies)
	})
}
