package brand

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	domain_brand "github.com/shodruzhoshimzoda/tojtech/internal/domain/brand"
)

type BrandRepository struct {
	dbPool *pgxpool.Pool
}

func NewBrandRepository(dbPool *pgxpool.Pool) *BrandRepository {
	return &BrandRepository{dbPool: dbPool}
}

// brandColumns - the same columns, in the same order, for every query that reads one brand.
// product_count is counted on the fly: it is not stored in the table.
const brandColumns = `
	b.uuid, b.name, b.slug, b.logo_url, b.description, b.created_at, b.updated_at,
	(SELECT COUNT(*) FROM products p WHERE p.brand_id = b.id)`

func scanBrand(row pgx.Row) (*domain_brand.Brand, error) {
	var b domain_brand.Brand
	err := row.Scan(&b.UUID, &b.Name, &b.Slug, &b.LogoURL, &b.Description, &b.CreatedAt, &b.UpdatedAt, &b.ProductCount)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// GetBrandByUUID returns one brand (product_count = all its products).
func (r *BrandRepository) GetBrandByUUID(ctx context.Context, id uuid.UUID) (*domain_brand.Brand, error) {
	query := `SELECT ` + brandColumns + ` FROM brands b WHERE b.uuid = $1`

	b, err := scanBrand(r.dbPool.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_brand.ErrBrandNotFound
		}
		return nil, fmt.Errorf("get brand: %w", err)
	}
	return b, nil
}

// CreateBrand inserts the brand and fills UUID, CreatedAt and UpdatedAt in the struct.
func (r *BrandRepository) CreateBrand(ctx context.Context, b *domain_brand.Brand) error {
	query := `
	INSERT INTO brands (name, slug, logo_url, description)
	VALUES ($1, $2, $3, $4)
	RETURNING uuid, created_at, updated_at
	`

	err := r.dbPool.QueryRow(ctx, query, b.Name, b.Slug, b.LogoURL, b.Description).
		Scan(&b.UUID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) { // name (ignoring case) or slug is taken
			return domain_brand.ErrBrandAlreadyExists
		}
		return fmt.Errorf("create brand: %w", err)
	}
	b.ProductCount = 0
	return nil
}

// UpdateBrand replaces name, slug, logo and description. updated_at is set by a trigger in the database.
func (r *BrandRepository) UpdateBrand(ctx context.Context, id uuid.UUID, b *domain_brand.Brand) (*domain_brand.Brand, error) {
	query := `
	UPDATE brands b
	SET name = $1, slug = $2, logo_url = $3, description = $4
	WHERE b.uuid = $5
	RETURNING ` + brandColumns

	updated, err := scanBrand(r.dbPool.QueryRow(ctx, query, b.Name, b.Slug, b.LogoURL, b.Description, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_brand.ErrBrandNotFound
		}
		if isUniqueViolation(err) {
			return nil, domain_brand.ErrBrandAlreadyExists
		}
		return nil, fmt.Errorf("update brand: %w", err)
	}
	return updated, nil
}

// DeleteBrand removes the brand. Its products stay: the database sets their brand_id to NULL
// (ON DELETE SET NULL in migration 000007).
func (r *BrandRepository) DeleteBrand(ctx context.Context, id uuid.UUID) error {
	tag, err := r.dbPool.Exec(ctx, `DELETE FROM brands WHERE uuid = $1`, id)
	if err != nil {
		return fmt.Errorf("delete brand: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain_brand.ErrBrandNotFound
	}
	return nil
}

// ListBrands is for the admin panel: ALL brands (even without products),
// product_count = all products of the brand, active or hidden.
// This is the number shown in the "delete this brand?" confirmation.
func (r *BrandRepository) GetBrands(ctx context.Context) ([]*domain_brand.Brand, error) {
	query := `SELECT ` + brandColumns + ` FROM brands b ORDER BY LOWER(b.name), b.id`

	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list brands: %w", err)
	}
	defer rows.Close()

	brands := make([]*domain_brand.Brand, 0)
	for rows.Next() {
		b, err := scanBrand(rows)
		if err != nil {
			return nil, fmt.Errorf("list brands: %w", err)
		}
		brands = append(brands, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list brands: %w", err)
	}
	return brands, nil
}

// ListBrandsWithProducts is for the shop's filter panel: only brands that have at least one
// ACTIVE product, and product_count = number of those active products.
// If categoryUUID is not nil, only products of that category are counted, so for "Graphics cards"
// the panel shows ASUS (12), MSI (7) and no HP if HP has no graphics cards.
func (r *BrandRepository) GetBrandsWithProducts(ctx context.Context, categoryUUID *uuid.UUID) ([]*domain_brand.Brand, error) {
	query := `
	SELECT b.uuid, b.name, b.slug, b.logo_url, b.description, b.created_at, b.updated_at, COUNT(p.id)
	FROM brands b
	JOIN products p ON p.brand_id = b.id AND p.is_active = TRUE
	WHERE $1::uuid IS NULL
	   OR p.category_id = (SELECT c.id FROM categories c WHERE c.uuid = $1::uuid)
	GROUP BY b.id
	ORDER BY LOWER(b.name), b.id
	`

	var category any // nil -> SQL NULL -> "no category filter"
	if categoryUUID != nil {
		category = *categoryUUID
	}

	rows, err := r.dbPool.Query(ctx, query, category)
	if err != nil {
		return nil, fmt.Errorf("list brands with products: %w", err)
	}
	defer rows.Close()

	brands := make([]*domain_brand.Brand, 0)
	for rows.Next() {
		b, err := scanBrand(rows)
		if err != nil {
			return nil, fmt.Errorf("list brands with products: %w", err)
		}
		brands = append(brands, b)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list brands with products: %w", err)
	}
	return brands, nil
}