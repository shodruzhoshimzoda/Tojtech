// pkg/reqlog/reqlog.go
package reqlog

import (
	"context"
	"log/slog"
	"sync"
)

type ctxKey struct{}

type entry struct {
	mu    sync.Mutex
	level slog.Level
	msg   string
	attrs []slog.Attr
}

// Inject - вызывается один раз в middleware, кладёт пустой накопитель в контекст
func Inject(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKey{}, &entry{level: slog.LevelInfo, msg: "request completed"})
}

// Info - хендлер помечает успешный исход
func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	set(ctx, slog.LevelInfo, msg, attrs...)
}

// Warn - хендлер помечает некритичную проблему (валидация, 404 и т.п.)
func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	set(ctx, slog.LevelWarn, msg, attrs...)
}

// Error - хендлер помечает реальную ошибку; уровень строки станет ERROR
func Error(ctx context.Context, msg string, err error, attrs ...slog.Attr) {
	attrs = append(attrs, slog.String("err", err.Error()))
	set(ctx, slog.LevelError, msg, attrs...)
}

func set(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	if e, ok := ctx.Value(ctxKey{}).(*entry); ok {
		e.mu.Lock()
		e.level = level
		e.msg = msg
		e.attrs = append(e.attrs, attrs...)
		e.mu.Unlock()
	}
}

// Extract - вызывается только из middleware в самом конце запроса
func Extract(ctx context.Context) (slog.Level, string, []slog.Attr) {
	if e, ok := ctx.Value(ctxKey{}).(*entry); ok {
		e.mu.Lock()
		defer e.mu.Unlock()
		return e.level, e.msg, e.attrs
	}
	return slog.LevelInfo, "request completed", nil
}
