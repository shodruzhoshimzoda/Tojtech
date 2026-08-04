package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	domain_category "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
	usecase "github.com/shodruzhoshimzoda/tojtech/internal/usecase/category"
	"github.com/shodruzhoshimzoda/tojtech/pkg/reqlog"
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
		reqlog.Warn(r.Context(), "category not found")
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, map[string]string{"error": "category not found"})
		return
	}

	categories, err := c.ucs.GetCategory(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain_category.ErrCategoryNotFound) {
			reqlog.Warn(r.Context(), "category not found")
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "category not found"})
			return
		}

		reqlog.Error(r.Context(), "failed to get category", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "internal server error"})
		return
	}

	render.JSON(w, r, map[string]any{"categories": categories})

}
