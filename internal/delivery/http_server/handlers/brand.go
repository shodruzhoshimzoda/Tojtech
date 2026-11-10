package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"

	domain_brand "github.com/shodruzhoshimzoda/tojtech/internal/domain/brand"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	brand_usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/brand"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

// A logo can be sent as an uploaded image (data URL), so the body may be big,
// but not unlimited. The exact logo limit is checked in the domain (MaxLogoDataURLLength).
const maxBrandBodyBytes = 4 << 20 // 4 MB

type BrandHandler struct {
	usc *brand_usecase.BrandUsecase
}

func NewBrandHandler(uc *brand_usecase.BrandUsecase) *BrandHandler {
	return &BrandHandler{usc: uc}
}

// GetShopBrands - GET /brands  (public)
// Brands for the filter panel: only brands that have active products.
// Optional ?category=<category uuid> counts only the products of that category.
func (h *BrandHandler) GetShopBrands(w http.ResponseWriter, r *http.Request) {
	const op = "BrandHandler.GetShopBrands"

	var categoryUUID *uuid.UUID
	if raw := r.URL.Query().Get("category"); raw != "" {
		id, err := uuid.Parse(raw)
		if err != nil {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid category uuid in query", "invalid category id")
			return
		}
		categoryUUID = &id
	}

	brands, err := h.usc.GetShopBrands(r.Context(), categoryUUID)
	if err != nil {
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to list shop brands", err, "internal server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"brands": brands})
}

// GetAllBrands - GET /brands/all  (admin)
// Every brand, even without products; product_count = all products of the brand.
func (h *BrandHandler) GetAllBrands(w http.ResponseWriter, r *http.Request) {
	const op = "BrandHandler.GetAllBrands"

	brands, err := h.usc.GetBrands(r.Context())
	if err != nil {
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to list brands", err, "internal server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"brands": brands})
}

// GetBrand - GET /brands/{uuid}  (public)
func (h *BrandHandler) GetBrand(w http.ResponseWriter, r *http.Request) {
	const op = "BrandHandler.GetBrand"

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid brand uuid", "invalid brand id")
		return
	}

	b, err := h.usc.GetBrand(r.Context(), id)
	if err != nil {
		h.respondBrandError(w, r, err, "failed to get brand", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"brand": b})
}

// CreateBrand - POST /brands  (admin)
func (h *BrandHandler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	const op = "BrandHandler.CreateBrand"

	var req dto.BrandRequest
	if !decodeBrandRequest(w, r, &req) {
		return
	}

	b, err := h.usc.CreateBrand(r.Context(), req)
	if err != nil {
		h.respondBrandError(w, r, err, "failed to create brand", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusCreated, map[string]any{"brand": b})
}

// UpdateBrand - PATCH /brands/{uuid}  (admin)
func (h *BrandHandler) UpdateBrand(w http.ResponseWriter, r *http.Request) {
	const op = "BrandHandler.UpdateBrand"

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid brand uuid", "invalid brand id")
		return
	}

	var req dto.BrandRequest
	if !decodeBrandRequest(w, r, &req) {
		return
	}

	b, err := h.usc.UpdateBrand(r.Context(), id, req)
	if err != nil {
		h.respondBrandError(w, r, err, "failed to update brand", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"brand": b})
}

// DeleteBrand - DELETE /brands/{uuid}  (admin)
// The products of the brand stay, they just lose the brand.
func (h *BrandHandler) DeleteBrand(w http.ResponseWriter, r *http.Request) {
	const op = "BrandHandler.DeleteBrand"

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "invalid brand uuid", "invalid brand id")
		return
	}

	if err := h.usc.DeleteBrand(r.Context(), id); err != nil {
		h.respondBrandError(w, r, err, "failed to delete brand", op)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// decodeBrandRequest reads the JSON body (not bigger than maxBrandBodyBytes).
// If it fails, the answer is already sent and false is returned.
func decodeBrandRequest(w http.ResponseWriter, r *http.Request, req *dto.BrandRequest) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxBrandBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to decode request body", "invalid request body")
		return false
	}
	return true
}

// respondBrandError turns a domain error into an HTTP answer.
func (h *BrandHandler) respondBrandError(w http.ResponseWriter, r *http.Request, err error, logMsg, op string) {
	switch {
	case isBrandValidationError(err):
		httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, logMsg, "invalid brand", err.Error())
	case errors.Is(err, domain_brand.ErrBrandAlreadyExists):
		httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusConflict, logMsg, "brand already exists", err.Error())
	case errors.Is(err, domain_brand.ErrBrandNotFound):
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "brand not found", "brand not found")
	default:
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, logMsg, err, "internal server error", op)
	}
}

func isBrandValidationError(err error) bool {
	for _, target := range []error{
		domain_brand.ErrEmptyBrandName,
		domain_brand.ErrBrandNameTooLong,
		domain_brand.ErrEmptySlug,
		domain_brand.ErrInvalidSlug,
		domain_brand.ErrLongDescription,
		domain_brand.ErrInvalidLogoURL,
		domain_brand.ErrLogoTooLarge,
	} {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}