package httpserver

import (
	"log/slog"
	"net/http"
	"time"

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
	Product  	*handler.ProductHandler
	Category 	*handler.CategoryHandler
	Auth     	*handler.AuthHandler
	AdminUser 	*handler.AdminUserHandler
	Brand 		*handlers.BrandHandler
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

		requireAuth := func(r chi.Router) {
		r.Use(mwlogger.RequireAuth(deps.JWTSecret))
	}

	router.Route("/api/v1", func(r chi.Router) {
		mountProfileRoutes(r, deps.Handlers.Auth, requireAuth)
		mountAuthRoutes(r, deps.Handlers.Auth)
		mountCategoryRoutes(r, deps.Handlers.Category, requireAdmin)
		mountProductRoutes(r, deps.Handlers.Product, requireAdmin)
		mountAdminUserRoutes(r, deps.Handlers.AdminUser, requireAdmin)
		mountBrandRoutes(r, deps.Handlers.Brand, requireAdmin)
	})

	

	return router
}

func mountAuthRoutes(r chi.Router, h *handler.AuthHandler) {
	// at most 10 login/register attempts per minute from one IP (both routes share the counter)
	authLimit := mwlogger.RateLimit(10, time.Minute)

	r.Route("/auth", func(r chi.Router) {
		r.Use(mwlogger.LimitBody(64 << 10)) // 64 KB is plenty for an email and a password
		r.With(authLimit).Post("/register", h.RegisterUser)
		r.With(authLimit).Post("/login", h.LoginUser)
		r.Post("/refresh", h.Refresh)
		r.Post("/logout", h.Logout)
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

// collect handler for get all users for admin
func mountAdminUserRoutes(r chi.Router, h *handler.AdminUserHandler, requireAdmin func(chi.Router)) {
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			requireAdmin(r)
			r.Get("/", h.GetUsers)
			r.Patch("/{uuid}/role", h.ChangeRole)
		})
	})
}


func mountBrandRoutes(r chi.Router, h *handler.BrandHandler, requireAdmin func(chi.Router)) {
	r.Route("/brands", func(r chi.Router) {
		r.Get("/", h.GetShopBrands) // for the filter panel of the shop
		r.Get("/{uuid}", h.GetBrand)

		r.Group(func(r chi.Router) {
			requireAdmin(r)
			r.Get("/all", h.GetAllBrands) // for the admin panel: all brands with all products counted
			r.Post("/", h.CreateBrand)
			r.Patch("/{uuid}", h.UpdateBrand)
			r.Delete("/{uuid}", h.DeleteBrand)
		})
	})
}

// the profile of the current user: any logged-in user, no role check
func mountProfileRoutes(r chi.Router, h *handler.AuthHandler, requireAuth func(chi.Router)) {
	// changing the password needs the current one, so it is a target for guessing: 5 per minute
	passwordLimit := mwlogger.RateLimit(5, time.Minute)

	r.Route("/me", func(r chi.Router) {
		requireAuth(r)
		r.Use(mwlogger.LimitBody(64 << 10))
		r.Get("/", h.GetMe)
		r.Patch("/", h.UpdateMe)
		r.With(passwordLimit).Post("/password", h.ChangePassword)
	})
}