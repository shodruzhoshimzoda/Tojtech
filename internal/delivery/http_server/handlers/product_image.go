package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	domain_product "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type addImageRequest struct {
	ImageURL string `json:"image_url"`
	IsMain   bool   `json:"is_main"`
}

func NewProductImageDTO(img *domain_product.ProductImage) dto.ProductImageDTO {
	return dto.ProductImageDTO{UUID: img.UUID, ImageURL: img.ImageURL, IsMain: img.IsMain}
}

func (h *ProductHandler) AddProductImageHandler(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.AddProductImageHandler"

	productUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid uuid", "invalid product id")
		return
	}

	var req addImageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to parse request body", "invalid request body")
		return
	}

	img, err := h.usc.AddProductImage(r.Context(), productUUID, req.ImageURL, req.IsMain)
	if err != nil {
		switch {
		case errors.Is(err, domain_product.ErrEmptyImageURL), errors.Is(err, domain_product.ErrInvalidImageURL):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid image url", err.Error())
		case errors.Is(err, domain_product.ErrTooManyImages):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "too many images", "product cannot have more than 10 images")
		case errors.Is(err, domain_product.ErrProductNotFound):
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "product not found", "product not found")
		default:
			httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to add image", err, "internal server error", op)
		}
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusCreated, map[string]any{"image": NewProductImageDTO(img)})
}

func (h *ProductHandler) DeleteProductImageHandler(w http.ResponseWriter, r *http.Request) {
	const op = "ProductHandler.DeleteProductImageHandler"

	productUUID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid uuid", "invalid product id")
		return
	}

	imageUUID, err := uuid.Parse(chi.URLParam(r, "image_uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid uuid", "invalid image id")
		return
	}

	if err := h.usc.DeleteProductImage(r.Context(), productUUID, imageUUID); err != nil {
		if errors.Is(err, domain_product.ErrImageNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "image not found", "image not found")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to delete image", err, "internal server error", op)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
