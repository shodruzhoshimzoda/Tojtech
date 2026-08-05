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
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/category"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type CategoryHandler struct {
	logger *slog.Logger
	ucs    *usecase.CategoryUseCase
}

func NewCategoryHandler(logger *slog.Logger, ucs *usecase.CategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		logger: logger,
		ucs:    ucs,
	}
}

func (c *CategoryHandler) GetCategoryHandler(w http.ResponseWriter, r *http.Request) {

	const op = "CategoryHandler.GetCategoryHandler"
	uuID := chi.URLParam(r, "uuid")

	id, err := uuid.Parse(uuID)
	if err != nil {

		httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "could not parse the uuid", "product not found")
		return
	}

	categories, err := c.ucs.GetCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain_category.ErrCategoryNotFound) {

			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "category not found", "product not found")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to get category", err, "internal server error", op)
		return

	}

	render.JSON(w, r, map[string]any{"categories": categories})

}
func (c *CategoryHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandler.CreateCategoryHandler"

	var req domain_category.Category

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to create category", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to validate category", "category is invalid", err.Error())
		return
	}

	uuID, err := c.ucs.CreateCategory(r.Context(), &req)
	if err != nil {
		if errors.Is(err, domain_category.ErrCategoryAlreadyExists) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "failed to create category", "category already exists")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create category", err, "internal server error", op)
		return
	}

	req.UUID = uuID
	httphelpers.RespondJSON(w, r, http.StatusCreated, map[string]any{"category": req})
}
