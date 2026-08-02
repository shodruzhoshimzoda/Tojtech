package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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

	id := chi.URLParam(r, "id")

	idInt, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		reqlog.Warn(r.Context(), "invalid product id", slog.String("op", op))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "invalid id"})
		return
	}

	product, err := h.usc.GetProduct(r.Context(), idInt)
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

	reqlog.Info(r.Context(), "product fetched", slog.String("op", op), slog.String("product_id", id))
	render.JSON(w, r, product)
}
