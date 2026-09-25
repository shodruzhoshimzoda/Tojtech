package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shodruzhoshimzoda/tojtech/internal/config"
	httpserver "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server"
	handler "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers"
	"github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres"
	categoryrepository "github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres/category"
	productrepository "github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres/product"
	userrepository "github.com/shodruzhoshimzoda/tojtech/internal/repository/postgres/user"
	categoryusecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/category"
	productusecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/product"
	userusecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/user"
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

	// our routes
	
	router := httpserver.NewRoutes(httpserver.RouterDeps{
		Handlers:   collectHandlers(db, cfg, log),
		Logger:     log,
		JWTSecret:  []byte(cfg.Jwt.Secret),
		CORSOrigin: cfg.HttpServer.CORSOrigin,
	})

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

// Collect all hadlers
func collectHandlers(db *pgxpool.Pool, cfg *config.Config, log *slog.Logger) httpserver.Handlers {
	productRepo := productrepository.NewProductRepository(db)
	productUC := productusecase.NewProductUsecase(productRepo)

	categoryRepo := categoryrepository.NewCategoryRepository(db)
	categoryUC := categoryusecase.NewCategoryUseCase(categoryRepo)

	userRepo := userrepository.NewUserRepository(db)
	userUC := userusecase.NewAuthUsecase(userRepo, []byte(cfg.Jwt.Secret), cfg.Jwt.TokenTTL)

	return httpserver.Handlers{
		Product:  handler.NewProductHandler(productUC, log),
		Category: handler.NewCategoryHandler(categoryUC, log),
		Auth:     handler.NewAuthHandler(userUC),
	}
}