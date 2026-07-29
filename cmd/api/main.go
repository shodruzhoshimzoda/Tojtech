package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/shodruzhoshimzoda/tojtech/internal/config"
	"github.com/shodruzhoshimzoda/tojtech/internal/delivery/handlers"
	middleware "github.com/shodruzhoshimzoda/tojtech/internal/delivery/handlers/middlwares"
	"github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres"
	repo "github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres/product"
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/product"
	"github.com/shodruzhoshimzoda/tojtech/pkg/logger"
)


func main() {

	
	cfg := config.MustLoadConfig()		// init configuration


	logger := logger.SetupLogger(cfg.Env)	// init logger

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

	repo := repo.NewProductRepository(db)

	service := usecase.NewProductUsecase(repo)

	handlers := handlers.NewProductHandler(service, logger)


	router := chi.NewRouter()

	// Подключаем наш новый красивый цветной логер
	router.Use(middleware.PrettyStructuredLogger())

	router.Get("/api/product/{id}", handlers.GetProductHandler)


	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Error starting server", slog.String("err", err.Error()))
		return
	}


	// TOOD: Start the HTTP server


}