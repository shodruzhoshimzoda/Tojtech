package http_server

import (
	"log/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers"
	mwlogger "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers/middlwares"
)

func NewRoutes(
	productHandler *handlers.ProductHandler,
	categoryHandler *handlers.CategoryHandler,
	log *slog.Logger,
) chi.Router {

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwlogger.RequestLogger(log)) // my custom logger
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/api/v1/products", productHandler.ListProductsHandler)
	router.Get("/api/v1/products/{uuid}", productHandler.GetProductHandler)
	router.Get("/api/v1/categories/{uuid}", categoryHandler.GetCategoryHandler)
	router.Post("/api/v1/categories", categoryHandler.CreateCategoryHandler)

	return router
}
