package http_server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers"
	mwlogger "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers/middlwares"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers" // замените на ваш пакет для JSON-ответов
)

func NewRoutes(
	productHandler *handlers.ProductHandler,
	categoryHandler *handlers.CategoryHandler,
	log *slog.Logger,
) chi.Router {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.RequestLogger(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusMethodNotAllowed, "method not allowed", "HTTP method is not supported for this route")
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "route not found", "the requested endpoint does not exist")
	})

	router.Route("/api/v1", func(r chi.Router) {

		// for categories
		r.Route("/categories", func(r chi.Router) {
			r.Get("/", categoryHandler.GetCategoriesHandler)
			r.Post("/", categoryHandler.CreateCategoryHandler)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", categoryHandler.GetCategoryHandler)
				r.Patch("/", categoryHandler.UpdateCategoryHandler)
				r.Put("/", categoryHandler.UpdateCategoryHandler)
				r.Delete("/", categoryHandler.DeleteCategoryHandler)
			})
		})

		// for products
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.ListProductsHandler)
			r.Post("/", productHandler.CreateProductHandler)

			r.Route("/{uuid}", func(r chi.Router) {
				r.Get("/", productHandler.GetProductHandler)
				r.Patch("/", productHandler.UpdateProductHandler)
				r.Put("/", productHandler.UpdateProductHandler)
				r.Delete("/", productHandler.DeleteProductHandler)
				r.Post("/images", productHandler.AddProductImageHandler)
				r.Delete("/images/{image_uuid}", productHandler.DeleteProductImageHandler)
			})
		})

	})

	return router
}
