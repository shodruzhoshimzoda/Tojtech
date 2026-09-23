package http_server

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers"
	mwlogger "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers/middlwares"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers" // замените на ваш пакет для JSON-ответов
)

func NewRoutes(
	productHandler *handlers.ProductHandler,
	categoryHandler *handlers.CategoryHandler,
	authHandler *handlers.AuthHandler,
	log *slog.Logger,
	jwtSecret []byte,
) chi.Router {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

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
			r.Get("/", categoryHandler.GetCategories)
			r.Get("/{uuid}", categoryHandler.GetCategory)

			r.Group(func(r chi.Router) {
				r.Use(mwlogger.RequireAuth(jwtSecret))
				r.Use(mwlogger.RequireRole(domain_user.RoleAdmin))

				r.Post("/", categoryHandler.CreateCategory)
				r.Patch("/{uuid}", categoryHandler.UpdateCategory)
				r.Delete("/{uuid}", categoryHandler.DeleteCategory)
			})

		})

		// for products
		r.Route("/products", func(r chi.Router) {
			r.Get("/", productHandler.GetProducts)
			r.Get("/{uuid}", productHandler.GetProduct)

			r.Group(func(r chi.Router) {
				r.Use(mwlogger.RequireAuth(jwtSecret))
				r.Use(mwlogger.RequireRole(domain_user.RoleAdmin))

				r.Post("/", productHandler.CreateProduct)
				r.Patch("/{uuid}", productHandler.UpdateProduct)
				r.Delete("/{uuid}", productHandler.DeleteProduct)
				r.Post("/{uuid}/images", productHandler.AddProductImageHandler)
				r.Delete("/{uuid}/images/{image_uuid}", productHandler.DeleteProductImageHandler)

			})
		})

		// for registration and authentication
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.RegisterUser)
			r.Post("/login", authHandler.LoginUser)
		})
	})

	return router
}

