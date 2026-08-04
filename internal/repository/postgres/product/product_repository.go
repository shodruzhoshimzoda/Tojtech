package repo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	domain_category "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
	domain_product "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
)

//

type ProductRepository struct {
	dbPool *pgxpool.Pool
}

func NewProductRepository(dbPool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		dbPool: dbPool,
	}
}

func (r *ProductRepository) GetByUUID(ctx context.Context, id uuid.UUID) (*domain_product.Product, error) {
	query := `
		SELECT 
			p.id, p.uuid, p.name, p.slug, p.description, p.price, p.stock, 
			p.category_id, p.is_active, p.created_at, p.updated_at,
			c.id, c.uuid, c.name, c.slug, c.description, c.created_at
		FROM products p
		INNER JOIN categories c ON p.category_id = c.id  -- вот JOIN
		WHERE p.uuid = $1 AND p.is_active = true
	`

	var product domain_product.Product
	var category domain_category.Category

	err := r.dbPool.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.UUID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.Stock,
		&product.CategoryID, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		&category.ID, &category.UUID, &category.Name, &category.Slug, &category.Description, &category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_product.ErrProductNotFound
		}
		return nil, err
	}

	product.Category = &category // кладем категорию внутрь продукта

	return &product, nil
}

func (r *ProductRepository) List(ctx context.Context) ([]domain_product.Product, error) {

	query := "SELECT  uuid, name, slug, description, price, stock, category_id, is_active, created_at, updated_at FROM products"

	rows, err := r.dbPool.Query(ctx, query)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []domain_product.Product{}, nil
		}
	}

	var p []domain_product.Product

	for rows.Next() {
		var product domain_product.Product
		err := rows.Scan(
			&product.UUID,
			&product.Name,
			&product.Slug,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CategoryID,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		p = append(p, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return p, nil
}
