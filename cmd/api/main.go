package main

import (
	"log/slog"

	"github.com/shodruzhoshimzoda/tojtech/internal/config"
	"github.com/shodruzhoshimzoda/tojtech/pkg/logger"
)


func main() {

	
	cfg := config.MustLoadConfig()
	
	logger := logger.SetupLogger(cfg.Env)

	logger.Info("Logger initialized", slog.String("env", cfg.Env))
	logger.Debug("DEBUG mode  enabled")



	// TOOD: Initialize database connection
	// TOOD: Initialize repository layer
	// TOOD: Initialize service (usecase) layer
	// TOOD: Initialize HTTP server and routes

	// TOOD: Start the HTTP server


}