package logger

import (
	"log/slog"
	"os"

	"github.com/shodruzhoshimzoda/tojtech/pkg/logger/handler/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// ANSI цвета

// SetupLogger - единая точка создания логгера приложения.
// local/dev -> красивый цветной вывод в консоль.
// prod      -> строгий структурированный JSON, без цвета.
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slogpretty.SetupPrettyLogger()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		// Если env передался некорректно (например, пустая строка или опечатка)
		// Фолбэчимся на дефолтный JSON или вывод с предупреждением
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		log.Warn("Unknown environment, fallback to default text logger", slog.String("env", env))
	}

	return log
}
