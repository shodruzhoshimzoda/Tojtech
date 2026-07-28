package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/shodruzhoshimzoda/tojtech/internal/config"
	"github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres"
	"github.com/shodruzhoshimzoda/tojtech/pkg/logger"
)


func main() {

	
	cfg := config.MustLoadConfig()
	
	logger := logger.SetupLogger(cfg.Env)

	logger.Info("Logger initialized", slog.String("env", cfg.Env))

	logger.Debug("DEBUG mode  enabled")

	ctx := context.Background()

	db, err  := postgres.ConnectionDB(ctx, cfg.GetDSN())

	if err != nil {
		logger.Error("unable connection to Database: ", slog.String("err",err.Error()))
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("Connection to database was successfully")

	
	// TOOD: Initialize repository layer
	// TOOD: Initialize service (usecase) layer
	// TOOD: Initialize HTTP server and routes

	// TOOD: Start the HTTP server


}