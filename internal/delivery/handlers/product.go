package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/product"
	"github.com/shodruzhoshimzoda/tojtech/pkg/reqlog"
)

type ProductHandler struct {
	usc *usecase.ProductUsecase
	log *slog.Logger
}

func NewProductHandler(usc *usecase.ProductUsecase, log *slog.Logger) *ProductHandler {
	return &ProductHandler{
		usc: usc,
		log: log,
	}
}

func (h *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {

	const op = "ProductHandler.GetProductHandler"

	uuID := chi.URLParam(r, "uuid")

	id, err := uuid.Parse(uuID)

	if err != nil {
		reqlog.Warn(r.Context(), "product not found", slog.String("op", op))
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "product not found"})
		return
	}

	product, err := h.usc.GetProduct(r.Context(), id)
	if err != nil {

		if errors.Is(err, domain.ErrProductNotFound) {
			reqlog.Warn(r.Context(), "product not found", slog.String("op", op))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "product not found"})
			return
		}

		reqlog.Error(r.Context(), "failed to get product", err, slog.String("op", op))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return

	}

	render.JSON(w, r, map[string]any{"product": product})
}

func (h *ProductHandler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.ListProductsHandler"

	products, err := h.usc.List(r.Context())
	if err != nil {
		reqlog.Error(r.Context(), "failed to list products", err, slog.String("op", op))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.JSON(w, r, map[string]any{"products": products})

}
