package repo_category

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	query := "SELECT uuid, name, slug, description, created_at, updated_at FROM categories WHERE uuid = $1"

	var category domain_category.Category

	err := c.dbPool.QueryRow(ctx, query, uuid).Scan(
		&category.UUID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_category.ErrCategoryNotFound
		}
		return nil, err
	}

	return &category, nil

}

// TODO: Implement other methods for the CategoryRepository as needed, such as CreateCategory, UpdateCategory, DeleteCategory, etc.

func (r *CategoryRepository) CreateCategory(ctx context.Context, c *domain_category.Category) (uuid.UUID, error) {
	query := `
	INSERT INTO categories (name, slug, description, created_at)
	VALUES ($1, $2, $3, $4)
	RETURNING uuid, created_at, updated_at
	`

	err := r.dbPool.QueryRow(ctx, query,
		c.Name,
		c.Slug,
		c.Description,
		time.Now(),
	).Scan(&c.UUID, &c.CreatedAt, &c.UpdatedAt) // сразу кладем в структуру

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return uuid.Nil, domain_category.ErrCategoryAlreadyExists
		}
		return uuid.Nil, err
	}
	return c.UUID, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, uuId uuid.UUID, cat *domain_category.Category) (*domain_category.Category, error) {
	query := `
	UPDATE categories
		SET name = $1, description = $2,slug = $3, updated_at = $4
		WHERE uuid = $5
		RETURNING uuid, name, description, slug, created_at, updated_at
	`
	var result domain_category.Category
	err := r.dbPool.QueryRow(ctx, query,
		cat.Name,
		cat.Description,
		cat.Slug,
		cat.UpdatedAt,
		uuId,
	).Scan(
		&result.UUID,
		&result.Name,
		&result.Description,
		&result.Slug,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		// 1. Проверяем, существует ли обновляемая категория
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_category.ErrCategoryNotFound
		}

		// 2. Проверяем, не нарушает ли новое имя/slug уникальность (SQLSTATE 23505)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain_category.ErrCategoryAlreadyExists
		}

		return nil, fmt.Errorf("update category: %w", err)
	}
	return &result, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, uuid uuid.UUID) error {
	query := "DELETE FROM categories WHERE uuid = $1"

	_, err := r.dbPool.Exec(ctx, query, uuid)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain_category.ErrCategoryNotFound
		}

		return fmt.Errorf("delete category: %w", err)

	}

	return nil
}
