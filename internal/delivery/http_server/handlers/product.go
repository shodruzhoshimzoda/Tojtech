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
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
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
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "invalid uuid", "product not found")
		return
	}

	product, err := h.usc.GetProduct(r.Context(), id)
	if err != nil {

		if errors.Is(err, domain.ErrProductNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "product not found", "product not found")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to get product", err, "internal server error", op)
		return

	}

	render.JSON(w, r, map[string]any{"product": product})
}

func (h *ProductHandler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.ListProductsHandler"

	products, err := h.usc.ProductList(r.Context())
	if err != nil {
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to list products", err, "internal server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"products": products})
}
