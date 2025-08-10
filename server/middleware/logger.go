package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-Id") == "" {
			r.Header.Set("X-Request-Id", uuid.New().String())
		}

		ctx := r.Context()

		start := time.Now()
		log := log.Logger.With().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("request_id", r.Header.Get("X-Request-ID")).
			Logger()

		ctx = log.WithContext(ctx)

		next.ServeHTTP(w, r.WithContext(ctx))

		duration := time.Since(start)
		log.Info().
			Str("duration", duration.String()).
			Msg("request completed")
	})
}
