// internal/delivery/handlers/middlewares/logger.go
package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

// RequestLogger - HTTP middleware для логирования запросов.
// Принимает уже готовый *slog.Logger - не создаёт собственный handler.
func RequestLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			level := slog.LevelInfo
			switch {
			case ww.Status() >= 500:
				level = slog.LevelError
			case ww.Status() >= 400:
				level = slog.LevelWarn
			}

			log.LogAttrs(r.Context(), level, "HTTP request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.String("duration", duration.String()),
			)
		})
	}
}
