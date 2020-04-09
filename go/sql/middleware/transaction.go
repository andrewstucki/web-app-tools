package sql

import (
	"net/http"

	"github.com/andrewstucki/web-app-tools/go/common"
	"github.com/andrewstucki/web-app-tools/go/sql/context"

	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// Transaction is a middleware that wraps a handler in a transaction and commits
// the transaction if the status code to return is in the 2xx-3xx range
func Transaction(db *sqlx.DB, renderer common.Renderer, logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tx, err := db.BeginTxx(ctx, nil)
			if err != nil {
				logger.Error().Err(err).Msg("error while starting transaction")
				renderer.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			txCtx := context.WithTransaction(ctx, tx)
			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(wrapped, r.Clone(txCtx))

			status := wrapped.Status()
			if 200 <= status && status < 400 {
				if err := tx.Commit(); err != nil {
					logger.Error().Err(err).Msg("error while starting transaction")
					renderer.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
				return
			}
			if err := tx.Rollback(); err != nil {
				logger.Error().Err(err).Msg("error while starting transaction")
				renderer.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			}
		})
	}
}
