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
)

func main() {

	cfg := config.MustLoadConfig() // init configuration

	log := logger.SetupLogger(cfg.Env) // init logger

	log.Info("Logger initialized", slog.String("env", cfg.Env))

	log.Debug("DEBUG mode  enabled")

	ctx := context.Background()

	db, err := postgres.ConnectionDB(ctx, cfg.GetDSN())

	if err != nil {
		log.Error("unable connection to Database: ", slog.String("err", err.Error()))
		os.Exit(1)
	}

	defer db.Close()

	log.Info("Connection to database was successfully")

	rep := repo.NewProductRepository(db)

	service := usecase.NewProductUsecase(rep)

	hand := handlers.NewProductHandler(service, log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.RequestLogger(log)) // my custom logger
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Get("/api/product/{id}", hand.GetProductHandler)

	// TODO: Start the HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port),
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	log.Info("starting server", slog.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server error", slog.String("err", err.Error()))
		os.Exit(1)
	}

}

//
//func setupPrettyLogger() *slog.Logger {
//	opts := slogpretty.PrettyHandlerOptions{&slog.HandlerOptions{
//		Level: slog.LevelDebug,
//	}}
//	handler := opts.NewPrettyHandler(os.Stdout)
//
//	return slog.New(handler)
//}
