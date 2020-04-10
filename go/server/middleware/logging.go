package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

func RequestLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			startTime := time.Now()
			defer func() {
				status := wrapped.Status()
				level := zerolog.InfoLevel
				switch {
				case status >= 500:
					level = zerolog.ErrorLevel
					break
				case status >= 400:
					level = zerolog.WarnLevel
				}

				fields := logger.
					WithLevel(level).
					Str("system", "http").
					Str("span.kind", "server").
					Str("http.host", r.URL.Host).
					Str("http.path", r.URL.Path).
					Str("http.query", r.URL.RawQuery).
					Str("http.method", r.Method).
					Str("http.start_time", startTime.Format(time.RFC3339)).
					Int("http.status", status).
					Int("http.bytes_written", wrapped.BytesWritten()).
					Dur("http.duration", time.Since(startTime))
				if d, ok := r.Context().Deadline(); ok {
					fields = fields.Str("http.request.deadline", d.Format(time.RFC3339))
				}
				fields.Msg("finished handling request")
			}()

			next.ServeHTTP(wrapped, r)
		})
	}
}
