package dto

import (
	"github.com/google/uuid"
	"time"

	category_domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
)

type CategoryDTO struct {
	UUID uuid.UUID `json:"uuid"`
}

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
		UUID:        uuid.UUID(c.UUID),
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		CreatedAt:   c.CreatedAt.Format(time.DateTime),
		UpdatedAt:   c.UpdatedAt.Format(time.DateTime),
	}
}
