package repo_product

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *ProductRepository) GetProductByUUID(ctx context.Context, id uuid.UUID) (*domain_product.Product, error) {
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

func (r *ProductRepository) ProductList(ctx context.Context) ([]*domain_product.Product, error) {

	query := "SELECT  uuid, name, slug, description, price, stock, category_id, is_active, created_at, updated_at FROM products"

	rows, err := r.dbPool.Query(ctx, query)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*domain_product.Product{}, nil
		}
		return nil, err
	}
	defer rows.Close()
	var p []*domain_product.Product

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

		p = append(p, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return p, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, p *domain_product.Product) error {
	if p.Category == nil {
		return domain_category.ErrCategoryNotFound
	}

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	queryProduct := `
       INSERT INTO products (name, slug, description, price, stock, category_id)
       SELECT $1, $2, $3, $4, $5, c.id
       FROM categories c
       WHERE c.uuid = $6
       RETURNING id, uuid, created_at, updated_at
    `
	err = tx.QueryRow(ctx, queryProduct,
		p.Name, p.Slug, p.Description, p.Price, p.Stock, p.Category.UUID,
	).Scan(&p.ID, &p.UUID, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain_category.ErrCategoryNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain_product.ErrProductAlreadyExists
		}
		return fmt.Errorf("create product: %w", err)
	}

	// 3. Вставляем картинки через tx
	if len(p.Images) > 0 {
		queryImage := `
          INSERT INTO product_images (product_id, image_url, is_main)
          VALUES ($1, $2, $3)
          RETURNING id, uuid, image_url, is_main
       `
		for i := range p.Images {
			img := &p.Images[i]
			img.ProductID = p.ID

			err := tx.QueryRow(ctx, queryImage, img.ProductID, img.ImageURL, img.IsMain).
				Scan(&img.ID, &img.UUID, &img.ImageURL, &img.IsMain)
			if err != nil {
				return fmt.Errorf("create product image: %w", err) // defer tx.Rollback отменит и создание товара!
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
