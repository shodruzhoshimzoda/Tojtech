package handlers

import (
	// "context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	// "github.com/google/uuid"
	category_domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
	// "github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/category"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type CategoryHandler struct {
	logger *slog.Logger
	ucs    *usecase.CategoryUseCase
}

func NewCategoryHandler(ucs *usecase.CategoryUseCase, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{
		logger: logger,
		ucs:    ucs,
	}
}

func (c *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {

	const op = "CategoryHandler.GetCategoryHandler"

	uuID := chi.URLParam(r, "uuid")
	category, err := c.ucs.GetCategory(r.Context(), uuID)

	if err != nil {
		if errors.Is(err, category_domain.ErrInvalidUUID) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "could not parse the uuid", "category not found")
			return
		}
		if errors.Is(err, category_domain.ErrCategoryNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "category not found", "category not found")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to get category", err, "internal server error", op)
		return
	}

	render.JSON(w, r, map[string]any{"category": category})

}

func (c *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	op := "categories.GetCategoriesHandler"

	categories, err := c.ucs.GetCategories(r.Context())
	if err != nil {
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "error getting categories", err, "internal server error", op)
		return
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"categories": categories})

}

func (c *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandler.CreateCategoryHandler"

	var req category_domain.Category

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to create category", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to validate category", "category is invalid", err.Error())
		return
	}

	category, err := c.ucs.CreateCategory(r.Context(), &req)
	if err != nil {
		if errors.Is(err, category_domain.ErrCategoryAlreadyExists) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "failed to create category", "category already exists")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create category", err, "internal server error", op)
		return
	}
	httphelpers.RespondJSON(w, r, http.StatusCreated, map[string]any{"category": category})
}

func (c *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandler.UpdateCategoryHandler"

	uuID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to parse the uuid", "category not found")
		return
	}

	var req category_domain.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to create category", "invalid request body")
		return
	}

	newCategory, err := c.ucs.UpdateCategory(r.Context(), uuID, &req)
	if err != nil {

		switch err {

		case category_domain.ErrInvalidCategoryName:
			httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to parse the category", "invalid request body", err.Error())
			return

		case category_domain.ErrEmptyCategoryName:
			httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to parse the category", "invalid request body", err.Error())
			return
		case category_domain.ErrSlugEmpty:
			httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to parse the category", "invalid request body", err.Error())
			return
		case category_domain.ErrLongDescription:
			httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to parse the category", "invalid request body", err.Error())
			return
		//////////////////

		case category_domain.ErrCategoryNotFound:
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "category not found", "category not found")
			return
		case category_domain.ErrCategoryAlreadyExists:
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "failed to create category", "duplicate name")
			return
		default:
			httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create category", err, "internal server error", op)
			return
		}

	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"category": newCategory})

}

func (c *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandler.DeleteCategoryHandler"

	id, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "could not parse the uuid", "product not found")
		return
	}

	if err := c.ucs.DeleteCategory(r.Context(), id); err != nil {
		if errors.Is(err, category_domain.ErrCategoryNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "category not found", "category not found")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to delete category", err, "internal server error", op)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
