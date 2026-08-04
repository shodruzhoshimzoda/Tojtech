package repo_category

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domain_category "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
)

type CategoryRepository struct {
	dbPool *pgxpool.Pool
}

func NewCategoryRepository(dbPool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{dbPool: dbPool}

}

func (c *CategoryRepository) GetCategoryByUUID(ctx context.Context, uuid uuid.UUID) (*domain_category.Category, error) {
	query := "SELECT uuid, name, slug, description, created_at FROM categories WHERE uuid = $1"

	var category domain_category.Category

	err := c.dbPool.QueryRow(ctx, query, uuid).Scan(
		&category.UUID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_category.ErrCategoryNotFound
		}
	}

	return &category, nil

}

// TODO: Implement other methods for the CategoryRepository as needed, such as CreateCategory, UpdateCategory, DeleteCategory, etc.
