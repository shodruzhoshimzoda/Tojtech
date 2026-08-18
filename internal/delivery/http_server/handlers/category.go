package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	category_domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/category"
	"github.com/shodruzhoshimzoda/tojtech/pkg/httphelpers"
)

type CategoryResponse struct {
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// NewCategoryResponse - DTO
func NewCategoryResponse(c *category_domain.Category) CategoryResponse {
	return CategoryResponse{
		UUID:        c.UUID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		CreatedAt:   c.CreatedAt.Format(time.DateTime),
		UpdatedAt:   c.UpdatedAt.Format(time.DateTime),
	}
}

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

	id, err := uuid.Parse(uuID)
	if err != nil {

		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "could not parse the uuid", "category not found")
		return
	}

	category, err := c.ucs.GetCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, category_domain.ErrCategoryNotFound) {

			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "category not found", "category not found")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to get category", err, "internal server error", op)
		return

	}
	catDTO := NewCategoryResponse(category)

	render.JSON(w, r, map[string]any{"category": catDTO})

}

func (c *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	op := "categories.GetCategoriesHandler"

	categories, err := c.ucs.GetCategories(r.Context())
	if err != nil {
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "error getting categories", err, "internal server error", op)
		return
	}

	var categoriesDTO = make([]CategoryResponse, 0)

	for _, cat := range categories {
		catDTO := NewCategoryResponse(cat)
		categoriesDTO = append(categoriesDTO, catDTO)
	}

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"categories": categoriesDTO})

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

	_, err := c.ucs.CreateCategory(r.Context(), &req)
	if err != nil {
		if errors.Is(err, category_domain.ErrCategoryAlreadyExists) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "failed to create category", "category already exists")
			return
		}
		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to create category", err, "internal server error", op)
		return
	}
	catDTO := NewCategoryResponse(&req)
	httphelpers.RespondJSON(w, r, http.StatusCreated, map[string]any{"category": catDTO})
}

func (c *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	const op = "CategoryHandler.UpdateCategoryHandler"

	uuID := chi.URLParam(r, "uuid")
	id, err := uuid.Parse(uuID)
	if err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "could not parse the uuid", "product not found")
		return
	}

	var req category_domain.Category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelpers.RespondWarn(r.Context(), w, r, http.StatusBadRequest, "failed to parse request body", "invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		httphelpers.RespondWarnWithDesc(r.Context(), w, r, http.StatusBadRequest, "failed to validate category", "category is invalid", err.Error())
		return
	}

	newCategory, err := c.ucs.UpdateCategory(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, category_domain.ErrCategoryNotFound) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusNotFound, "category not found", "category not found")
			return
		}
		if errors.Is(err, category_domain.ErrCategoryAlreadyExists) {
			httphelpers.RespondWarn(r.Context(), w, r, http.StatusConflict, "category with this name already exists", "duplicate name")
			return
		}

		httphelpers.RespondError(r.Context(), w, r, http.StatusInternalServerError, "failed to update category", err, "internal server error", op)
		return
	}

	catDTO := NewCategoryResponse(newCategory)

	httphelpers.RespondJSON(w, r, http.StatusOK, map[string]any{"category": catDTO})

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
