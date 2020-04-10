package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/oauth"
	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
	sqlContext "github.com/andrewstucki/web-app-tools/go/sql/context"
)

var (
	currentUserKey = "current-user-context-key"
)

// CurrentUser gets the current user or errors
func CurrentUser(ctx context.Context) (interface{}, error) {
	ctxFn := ctx.Value(&currentUserKey)
	if ctxFn == nil {
		return nil, errors.New("the callback must be injected")
	}
	return ctxFn.(func(ctx context.Context) (interface{}, error))(ctx)
}

func setCurrentUserFn(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) context.Context {
	return context.WithValue(ctx, &currentUserKey, fn)
}

// currentUser is a middleware that injects the current user into the context
func currentUser(db *sqlx.DB, handler *oauth.Handler, renderer common.Renderer, logger zerolog.Logger, getter func(ctx context.Context, queryer sqlContext.QueryContext, claims *verifier.StandardClaims) (interface{}, error)) func(next http.Handler) http.Handler {
	fn := func(ctx context.Context) (interface{}, error) {
		claims := handler.Claims(ctx)
		if claims != nil {
			return getter(ctx, sqlContext.GetQueryer(ctx, db), claims)
		}
		return nil, nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.Clone(setCurrentUserFn(r.Context(), fn)))
		})
	}
}
