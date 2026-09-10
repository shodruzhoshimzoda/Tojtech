package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	domain_category "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	domain_product "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
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

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {

	const op = "ProductHandler.GetProductHandler"

	uuID := chi.URLParam(r, "uuid")

	id, err := uuid.Parse(uuID)

	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid uuid", "product not found")
		return
	}

	product, err := h.usc.GetProduct(r.Context(), id)
	if err != nil {

		if errors.Is(err, domain_product.ErrProductNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "product not found", "product not found")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to get product", err, "internal server error", op)
		return

	}
	prodDTO := dto.NewProductDTO(product)
	render.JSON(w, r, map[string]any{"product": prodDTO})
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.ListProductsHandler"

	products, err := h.usc.ProductList(r.Context())
	if err != nil {
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to list products", err, "internal server error", op)
		return
	}
	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"products": products})
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.CreateProductHandler"

	var prod domain_product.Product
	if err := json.NewDecoder(r.Body).Decode(&prod); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to create product", "invalid request body")
		return
	}

	if err := prod.Validate(); err != nil {
		httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to validate product", "product is invalid", err.Error())
		return
	}

	err := h.usc.CreateProduct(r.Context(), &prod)
	if err != nil {
		switch {
		case errors.Is(err, domain_category.ErrCategoryNotFound):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to create product", "specified category does not exist")
			return
		case errors.Is(err, domain_product.ErrProductAlreadyExists):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "failed to create product", "product with this slug already exists")
			return
		default:
			httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create product", err, "internal server error", op)
			return
		}
	}

	prodDTO := dto.NewProductDTO(&prod)
	httphelpers.RespondJSON(w, r, http.StatusCreated, map[string]any{"product": prodDTO})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.DeleteProductHandler"

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid uuid", "product not found")
		return
	}

	err = h.usc.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain_product.ErrProductNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "product not found", "product not found")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to delete product", err, "internal server error", op)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.UpdateProductHandler"

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid uuid", "invalid product id")
		return
	}

	var req domain_product.Product
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to parse request body", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to validate product", "product is invalid", err.Error())
		return
	}

	req.UUID = id // UUID из URL всегда главнее того, что могло прийти в теле запроса

	if err := h.usc.UpdateProduct(r.Context(), &req); err != nil {
		switch {
		case errors.Is(err, domain_product.ErrProductNotFound):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "product not found", "product not found")

		case errors.Is(err, domain_category.ErrCategoryNotFound):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to update product", "specified category does not exist")

		case errors.Is(err, domain_product.ErrProductAlreadyExists):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "failed to update product", "product with this slug already exists")

		default:
			httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to update product", err, "internal server error", op)
		}
		return
	}

	prodDTO := dto.NewProductDTO(&req)
	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"product": prodDTO})
}
