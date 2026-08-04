package mwlogger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/shodruzhoshimzoda/tojtech/pkg/reqlog"
)

func RequestLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := reqlog.Inject(r.Context())
			r = r.WithContext(ctx)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			level, msg, handlerAttrs := reqlog.Extract(r.Context())

			// статус 5xx/4xx всегда поднимает уровень, даже если хендлер сам этого не сделал
			switch {
			case ww.Status() >= 500 && level < slog.LevelError:
				level = slog.LevelError
			case ww.Status() >= 400 && level < slog.LevelWarn:
				level = slog.LevelWarn
			}

			attrs := []slog.Attr{
				slog.String("component", "middleware/logger"),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				//slog.String("remote_addr", r.RemoteAddr),
				//slog.String("user_agent", r.UserAgent()),
				//slog.String("request_id", middleware.GetReqID(r.Context())),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.String("duration", duration.String()),
			}
			attrs = append(attrs, handlerAttrs...)

			log.LogAttrs(r.Context(), level, msg, attrs...)
		})
	}
}
