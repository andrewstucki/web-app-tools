package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/rs/zerolog"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a request ID if one is provided.
func Recoverer(renderer common.Renderer, logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if caught := recover(); caught != nil && caught != http.ErrAbortHandler {
					logger.Error().Interface("panic", caught).Str("trace", string(debug.Stack())).Msg("recovering from panic")
					renderer.Render(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
