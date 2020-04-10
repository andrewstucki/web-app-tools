package server

import (
	"context"
	"net/http"

	"github.com/andrewstucki/web-app-tools/go/oauth"
)

var (
	tokenUserKey = "token-user-context-key"
)

func setToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, &tokenUserKey, token)
}

func getToken(ctx context.Context) string {
	token := ctx.Value(&tokenUserKey)
	if token == nil {
		return ""
	}
	return token.(string)
}

// tokenUser is a middleware that checks for an API token into the context
func tokenUser(handler *oauth.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if claims := handler.Claims(r.Context()); claims != nil {
				// Skip the check for an API token since we're already authed
				next.ServeHTTP(w, r)
				return
			}
			token := r.Header.Get("X-Api-Token")
			if token == "" {
				next.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r.Clone(setToken(r.Context(), token)))
		})
	}
}
