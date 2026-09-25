package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	handler "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers"
	handlers "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers"
	mwlogger "github.com/shodruzhoshimzoda/tojtech/internal/delivery/http_server/handlers/middlwares"
	userdomain "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type Handlers struct {
	Product  *handler.ProductHandler
	Category *handler.CategoryHandler
	Auth     *handler.AuthHandler
}

type RouterDeps struct {
	Handlers   Handlers
	Logger     *slog.Logger
	JWTSecret  []byte
	CORSOrigin string
}

func NewRoutes(deps RouterDeps) chi.Router {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{deps.CORSOrigin},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.Use(middleware.RequestID)
	router.Use(mwlogger.RequestLogger(deps.Logger))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusMethodNotAllowed, "method not allowed", "HTTP method is not supported for this route")
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "route not found", "the requested endpoint does not exist")
	})


	requireAdmin := func(r chi.Router) {
			r.Use(mwlogger.RequireAuth(deps.JWTSecret))
			r.Use(mwlogger.RequireRole(userdomain.RoleAdmin))

	}

	router.Route("/api/v1", func(r chi.Router) {
		mountAuthRoutes(r, deps.Handlers.Auth)
		mountCategoryRoutes(r, deps.Handlers.Category, requireAdmin)
		mountProductRoutes(r, deps.Handlers.Product, requireAdmin)
	})

	

	return router
}

func mountAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.RegisterUser)
		r.Post("/login", h.LoginUser)
	})
}

func mountCategoryRoutes(r chi.Router, h *handlers.CategoryHandler, requireAdmin func(chi.Router)) {
	r.Route("/categories", func(r chi.Router) {
		r.Get("/", h.GetCategories)
		r.Get("/{uuid}", h.GetCategory)

		r.Group(func(r chi.Router) {
			requireAdmin(r)
			r.Post("/", h.CreateCategory)
			r.Patch("/{uuid}", h.UpdateCategory)
			r.Delete("/{uuid}", h.DeleteCategory)
		})
	})
}

func mountProductRoutes(r chi.Router, h *handlers.ProductHandler, requireAdmin func(chi.Router)) {
	r.Route("/products", func(r chi.Router) {
		r.Get("/", h.GetProducts)
		r.Get("/{uuid}", h.GetProduct)

		r.Group(func(r chi.Router) {
			requireAdmin(r)
			r.Post("/", h.CreateProduct)
			r.Patch("/{uuid}", h.UpdateProduct)
			r.Delete("/{uuid}", h.DeleteProduct)
			r.Post("/{uuid}/images", h.AddProductImageHandler)
			r.Delete("/{uuid}/images/{image_uuid}", h.DeleteProductImageHandler)
		})
	})
}

