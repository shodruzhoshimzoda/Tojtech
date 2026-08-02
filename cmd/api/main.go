package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shodruzhoshimzoda/tojtech/internal/config"
	"github.com/shodruzhoshimzoda/tojtech/internal/delivery/handlers"
	mwlogger "github.com/shodruzhoshimzoda/tojtech/internal/delivery/handlers/middlwares"
	"github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres"
	repo "github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres/product"
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/product"
	"github.com/shodruzhoshimzoda/tojtech/pkg/logger"
	"github.com/shodruzhoshimzoda/tojtech/pkg/logger/handler/slogpretty"
)

func main() {

	cfg := config.MustLoadConfig() // init configuration

	logger := logger.SetupLogger(cfg.Env) // init logger

	logger.Info("Logger initialized", slog.String("env", cfg.Env))

	logger.Debug("DEBUG mode  enabled")

	ctx := context.Background()

	db, err := postgres.ConnectionDB(ctx, cfg.GetDSN())

	if err != nil {
		logger.Error("unable connection to Database: ", slog.String("err", err.Error()))
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("Connection to database was successfully")

	repo := repo.NewProductRepository(db)

	service := usecase.NewProductUsecase(repo)

	handlers := handlers.NewProductHandler(service, logger)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Get("/api/product/{id}", handlers.GetProductHandler)

	// TOOD: Start the HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	logger.Info("starting server", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", slog.String("err", err.Error()))
		os.Exit(1)
	}

}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{&slog.HandlerOptions{
		Level: slog.LevelDebug,
	}}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
