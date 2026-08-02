package slogpretty

import (
	"context"
	"fmt"
	"io"
	stdLog "log"
	"log/slog"
	"os"
	"strings"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts PrettyHandlerOptions
	slog.Handler
	l     *stdLog.Logger
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		opts:    opts,
		Handler: slog.NewTextHandler(out, opts.SlogOpts), // Используем TextHandler!
		l:       stdLog.New(out, "", 0),
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	// Собираем атрибуты
	var attrsBuilder strings.Builder

	// Добавляем сохраненные атрибуты хэндлера
	for _, a := range h.attrs {
		attrsBuilder.WriteString(fmt.Sprintf("%s=%v ", color.WhiteString(a.Key), a.Value.Any()))
	}

	// Добавляем атрибуты конкретной записи
	r.Attrs(func(a slog.Attr) bool {
		attrsBuilder.WriteString(fmt.Sprintf("%s=%v ", color.WhiteString(a.Key), a.Value.Any()))
		return true
	})

	timeStr := r.Time.Format("[15:04:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(
		timeStr,
		level,
		msg,
		strings.TrimSpace(attrsBuilder.String()),
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	newAttrs := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)

	return &PrettyHandler{
		opts:    h.opts,
		Handler: h.Handler,
		l:       h.l,
		attrs:   newAttrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		opts:    h.opts,
		Handler: h.Handler.WithGroup(name),
		l:       h.l,
		attrs:   h.attrs,
	}
}

func SetupPrettyLogger() *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
