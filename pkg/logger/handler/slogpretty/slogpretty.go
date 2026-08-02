package slogpretty

import (
	"context"
	"fmt"
	"io"
	stdLog "log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

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
		Handler: slog.NewTextHandler(out, opts.SlogOpts),
		l:       stdLog.New(out, "", 0),
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	levelText := r.Level.String() + ":"

	var level string
	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(levelText)
	case slog.LevelInfo:
		level = color.BlueString(levelText)
	case slog.LevelWarn:
		level = color.YellowString(levelText)
	case slog.LevelError:
		level = color.RedString(levelText)
	}

	var attrsBuilder strings.Builder

	for _, a := range h.attrs {
		attrsBuilder.WriteString(formatAttr(a))
	}
	r.Attrs(func(a slog.Attr) bool {
		attrsBuilder.WriteString(formatAttr(a))
		return true
	})

	timeStr := r.Time.Format("[15:04:05.000]")

	// сообщение тоже реагирует на уровень - ошибки сразу бросаются в глаза
	var msg string
	if r.Level >= slog.LevelError {
		msg = color.New(color.FgRed, color.Bold).Sprint(r.Message)
	} else {
		msg = color.CyanString(r.Message)
	}

	h.l.Println(
		timeStr,
		level,
		msg,
		strings.TrimSpace(attrsBuilder.String()),
	)

	return nil
}

// formatAttr - красит key=value по смыслу конкретного ключа
func formatAttr(a slog.Attr) string {
	key := color.WhiteString(a.Key)
	value := colorizeValue(a)
	return fmt.Sprintf("%s=%v ", key, value)
}

func colorizeValue(a slog.Attr) string {
	val := fmt.Sprintf("%v", a.Value.Any())

	switch a.Key {
	case "method":
		return color.New(color.FgCyan, color.Bold).Sprint(val)

	case "status":
		status, _ := strconv.Atoi(val)
		return colorForStatus(status).Sprint(val)

	case "err", "error":
		return color.New(color.FgRed, color.Bold).Sprint(val)

	case "path":
		return color.BlueString(val)

	case "duration":
		return colorForDuration(val).Sprint(val)

	default:
		return val
	}
}

func colorForStatus(status int) *color.Color {
	switch {
	case status >= 500:
		return color.New(color.FgRed)
	case status >= 400:
		return color.New(color.FgYellow)
	case status >= 200:
		return color.New(color.FgGreen)
	default:
		return color.New(color.FgWhite)
	}
}

func colorForDuration(d string) *color.Color {
	dur, err := time.ParseDuration(d)
	if err != nil {
		return color.New(color.FgWhite)
	}
	switch {
	case dur > 500*time.Millisecond:
		return color.New(color.FgRed)
	case dur > 100*time.Millisecond:
		return color.New(color.FgYellow)
	default:
		return color.New(color.FgWhite)
	}
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
