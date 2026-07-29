package middleware

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/middleware"
)

// ANSI Цветовые коды
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
)

// PrettyHandler реализует slog.Handler с ANSI подсвечиванием
type PrettyHandler struct {
	out   io.Writer
	mu    *sync.Mutex
	attrs []slog.Attr
}

func NewPrettyHandler(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		out: out,
		mu:  &sync.Mutex{},
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// 1. Форматирование уровня логирования (Level)
	var levelStr string
	switch r.Level {
	case slog.LevelDebug:
		levelStr = fmt.Sprintf("%s[DEBUG]%s", colorPurple, colorReset)
	case slog.LevelWarn:
		levelStr = fmt.Sprintf("%s[WARN ]%s", colorYellow, colorReset)
	case slog.LevelError:
		levelStr = fmt.Sprintf("%s[ERROR]%s", colorRed, colorReset)
	default:
		levelStr = fmt.Sprintf("%s[INFO ]%s", colorGreen, colorReset)
	}

	// Время
	timeStr := fmt.Sprintf("%s%s%s", colorGray, r.Time.Format("15:04:05.000"), colorReset)

	// Собираем атрибуты (из r.Attrs и из h.attrs)
	var method, path, statusStr, durationStr string

	fn := func(a slog.Attr) bool {
		switch a.Key {
		case "method":
			method = fmt.Sprintf("%s%-6s%s", colorCyan, a.Value.String(), colorReset)
		case "path":
			path = fmt.Sprintf("%s%s%s", colorBlue, a.Value.String(), colorReset)
		case "status":
			status := int(a.Value.Int64())
			statusColor := colorGreen
			if status >= 400 && status < 500 {
				statusColor = colorYellow
			} else if status >= 500 {
				statusColor = colorRed
			}
			statusStr = fmt.Sprintf("%s%d%s", statusColor, status, colorReset)
		case "duration":
			durationStr = fmt.Sprintf("%s%s%s", colorGray, a.Value.String(), colorReset)
		}
		return true
	}

	// Сначала обрабатываем прикрепленные атрибуты хэндлера, затем атрибуты записи
	for _, attr := range h.attrs {
		fn(attr)
	}
	r.Attrs(fn)

	// Печать красивой строки в консоль
	if method != "" {
		fmt.Fprintf(h.out, "%s %s | %s | %s | %s | %s\n",
			timeStr, levelStr, statusStr, durationStr, method, path)
	} else {
		fmt.Fprintf(h.out, "%s %s %s\n", timeStr, levelStr, r.Message)
	}

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		out:   h.out,
		mu:    h.mu,
		attrs: append(h.attrs, attrs...),
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return h
}

func PrettyStructuredLogger() func(next http.Handler) http.Handler {
	handler := NewPrettyHandler(os.Stdout)
	log := slog.New(handler)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			next.ServeHTTP(ww, r)

			duration := time.Since(start)

			var level slog.Level
			switch {
			case ww.Status() >= 500:
				level = slog.LevelError
			case ww.Status() >= 400:
				level = slog.LevelWarn
			default:
				level = slog.LevelInfo
			}

			log.Log(r.Context(), level, "HTTP Request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.String("duration", duration.String()),
			)
		})
	}
}