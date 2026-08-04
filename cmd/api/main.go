package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/shodruzhoshimzoda/tojtech/internal/config"
	"github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server"
	"github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers"
	"github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres"
	repo_category "github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres/category"
	"github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres/product"
	usecase_category "github.com/shodruzhoshimzoda/tojtech/internal/usecase/category"
	"github.com/shodruzhoshimzoda/tojtech/internal/usecase/product"
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

	// for products
	repProd := repo_product.NewProductRepository(db)
	serviceProd := usecase_product.NewProductUsecase(repProd)
	handProd := handlers.NewProductHandler(serviceProd, log)

	// for repositories
	repCateg := repo_category.NewCategoryRepository(db)
	serCateg := usecase_category.NewCategoryUseCase(repCateg)
	handCateg := handlers.NewCategoryHandler(log, serCateg)

	// our routes
	router := http_server.NewRoutes(handProd,handCateg,  log)

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
