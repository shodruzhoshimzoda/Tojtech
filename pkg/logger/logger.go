// pkg/logger/logger.go
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// ANSI цвета
const (
	colorReset  = "\033[0m"
	colorGray   = "\033[90m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// SetupLogger - единая точка создания логгера приложения.
// local/dev -> красивый цветной вывод в консоль.
// prod      -> строгий структурированный JSON, без цвета.
func SetupLogger(env string) *slog.Logger {
	switch env {
	case envProd:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true, // в проде полезно знать файл:строку при ошибке
		}))
	case envDev:
		return slog.New(NewPrettyHandler(os.Stdout, slog.LevelDebug))
	default: // envLocal и любое неизвестное значение - безопасный дефолт
		return slog.New(NewPrettyHandler(os.Stdout, slog.LevelDebug))
	}
}

// ---------- PrettyHandler ----------

type PrettyHandler struct {
	out       io.Writer
	minLevel  slog.Level
	mu        *sync.Mutex
	attrs     []slog.Attr
	groupName string
}

func NewPrettyHandler(out io.Writer, minLevel slog.Level) *PrettyHandler {
	return &PrettyHandler{
		out:      out,
		minLevel: minLevel,
		mu:       &sync.Mutex{},
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.minLevel
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	levelStr, icon, color := levelStyle(r.Level)

	var b strings.Builder

	b.WriteString(colorGray)
	b.WriteString(r.Time.Format("15:04:05.000"))
	b.WriteString(colorReset)
	b.WriteByte(' ')

	b.WriteString(color)
	b.WriteString(icon)
	fmt.Fprintf(&b, "%-6s", levelStr) // выравнивание всё равно нужно через Fprintf
	b.WriteString(colorReset)
	b.WriteByte(' ')

	b.WriteString(colorBold)
	b.WriteString(r.Message)
	b.WriteString(colorReset)

	for _, a := range h.attrs {
		writeAttr(&b, a, h.groupName)
	}
	r.Attrs(func(a slog.Attr) bool {
		writeAttr(&b, a, h.groupName)
		return true
	})

	_, err := fmt.Fprintln(h.out, b.String())
	return err
}

func writeAttr(b *strings.Builder, a slog.Attr, group string) {
	key := a.Key
	if group != "" {
		key = group + "." + key
	}
	fmt.Fprintf(b, "  %s%s%s=%s", colorCyan, key, colorReset, formatValue(a.Value))
}

func formatValue(v slog.Value) string {
	s := v.String()
	if strings.ContainsAny(s, " \t") {
		return fmt.Sprintf("%q", s) // "два слова" -> в кавычках, чтобы не ломать парсинг глазами
	}
	return s
}

func levelStyle(level slog.Level) (text, icon, color string) {
	switch {
	case level >= slog.LevelError:
		return "ERROR", "✖ ", colorRed
	case level >= slog.LevelWarn:
		return "WARN", "▲ ", colorYellow
	case level >= slog.LevelInfo:
		return "INFO", "● ", colorGreen
	default:
		return "DEBUG", "○ ", colorPurple
	}
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)
	return &PrettyHandler{
		out: h.out, minLevel: h.minLevel, mu: h.mu,
		attrs: newAttrs, groupName: h.groupName,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	newGroup := name
	if h.groupName != "" {
		newGroup = h.groupName + "." + name
	}
	return &PrettyHandler{
		out: h.out, minLevel: h.minLevel, mu: h.mu,
		attrs: h.attrs, groupName: newGroup,
	}
}

var _ = time.Now // подавляем неиспользуемый импорт, если понадобится позже
