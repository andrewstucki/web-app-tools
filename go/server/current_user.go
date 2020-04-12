package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/rs/zerolog"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/oauth"
	"github.com/andrewstucki/web-app-tools/go/oauth/verifier"
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

// SetCurrentUserFn sets the context with the given current user resolver, it shouldn't
// be used directly and is only exported for testing
func SetCurrentUserFn(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) context.Context {
	return context.WithValue(ctx, &currentUserKey, fn)
}

// ClaimsOrToken represents either claims found
// or and API token found
type ClaimsOrToken struct {
	Claims *verifier.StandardClaims
	Token  string
}

// currentUser is a middleware that injects the current user into the context
func currentUser(handler *oauth.Handler, renderer common.Renderer, logger zerolog.Logger, getter func(ctx context.Context, claims *ClaimsOrToken) (interface{}, error)) func(next http.Handler) http.Handler {
	fn := func(ctx context.Context) (interface{}, error) {
		if claims := handler.Claims(ctx); claims != nil {
			return getter(ctx, &ClaimsOrToken{Claims: claims})
		}
		if token := getToken(ctx); token != "" {
			return getter(ctx, &ClaimsOrToken{Token: token})
		}
		return nil, nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.Clone(SetCurrentUserFn(r.Context(), fn)))
		})
	}
}
