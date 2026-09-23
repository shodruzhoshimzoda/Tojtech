package product_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	category_domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
	product_domain "github.com/shodruzhoshimzoda/tojtech/internal/domain/product"
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
func (r *ProductRepository) GetProduct(ctx context.Context, id uuid.UUID) (*product_domain.Product, error) {
	query := `
       SELECT 
          p.id, p.uuid, p.name, p.slug, p.description, p.price, p.stock, 
          p.category_id, p.is_active, p.created_at, p.updated_at,
          c.id, c.uuid, c.name, c.slug, c.description, c.created_at
       FROM products p
       INNER JOIN categories c ON p.category_id = c.id
       WHERE p.uuid = $1 AND p.is_active = true
    `

	var product product_domain.Product
	var category category_domain.Category

	err := r.dbPool.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.UUID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.Stock,
		&product.CategoryID, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		&category.ID, &category.UUID, &category.Name, &category.Slug, &category.Description, &category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_domain.ErrProductNotFound
		}
		return nil, err
	}

	product.Category = &category

	// --- ДОБАВЛЯЕМ ПОДГРУЗКУ КАРТИНОК ---
	queryImages := `
       SELECT uuid, image_url, is_main 
       FROM product_images 
       WHERE product_id = $1
    `
	rows, err := r.dbPool.Query(ctx, queryImages, product.ID)
	if err != nil {
		return nil, fmt.Errorf("get product images: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var img product_domain.ProductImage
		if err := rows.Scan(&img.UUID, &img.ImageURL, &img.IsMain); err != nil {
			return nil, fmt.Errorf("scan product image: %w", err)
		}
		product.Images = append(product.Images, img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context) ([]*product_domain.Product, error) {
	query := `
		SELECT
			p.id, p.uuid, p.name, p.slug, p.description, p.price, p.stock,
			p.category_id, p.is_active, p.created_at, p.updated_at,
			c.id, c.uuid, c.name, c.slug, c.description, c.created_at
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.is_active = true
	`

	rows, err := r.dbPool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*product_domain.Product, 0)
	productIDs := make([]int64, 0)

	for rows.Next() {
		var product product_domain.Product

		// LEFT JOIN - у товара может не быть категории, поля тогда NULL.
		// pgx умеет писать NULL прямо в nil-указатель, если сканировать
		// в **T (адрес указателя), а не в *T напрямую.
		var (
			catID        *int64
			catUUID      *uuid.UUID
			catName      *string
			catSlug      *string
			catDesc      *string
			catCreatedAt *time.Time
		)

		err := rows.Scan(
			&product.ID, &product.UUID, &product.Name, &product.Slug, &product.Description, &product.Price, &product.Stock,
			&product.CategoryID, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
			&catID, &catUUID, &catName, &catSlug, &catDesc, &catCreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if catID != nil {
			category := category_domain.Category{
				ID:        *catID,
				UUID:      *catUUID,
				Name:      *catName,
				Slug:      *catSlug,
				CreatedAt: *catCreatedAt,
			}
			if catDesc != nil {
				category.Description = *catDesc
			}
			product.Category = &category
		}

		products = append(products, &product)
		productIDs = append(productIDs, product.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(productIDs) == 0 {
		return products, nil
	}

	// Одним запросом забираем картинки сразу для всех товаров списка,
	// вместо отдельного запроса на каждый товар (N+1).
	imgRows, err := r.dbPool.Query(ctx, `
		SELECT product_id, uuid, image_url, is_main
		FROM product_images
		WHERE product_id = ANY($1)
	`, productIDs)
	if err != nil {
		return nil, fmt.Errorf("get product images: %w", err)
	}
	defer imgRows.Close()

	imagesByProduct := make(map[int64][]product_domain.ProductImage)
	for imgRows.Next() {
		var productID int64
		var img product_domain.ProductImage
		if err := imgRows.Scan(&productID, &img.UUID, &img.ImageURL, &img.IsMain); err != nil {
			return nil, fmt.Errorf("scan product image: %w", err)
		}
		imagesByProduct[productID] = append(imagesByProduct[productID], img)
	}
	if err := imgRows.Err(); err != nil {
		return nil, err
	}

	for _, p := range products {
		p.Images = imagesByProduct[p.ID]
	}

	return products, nil
}
func (r *ProductRepository) CreateProduct(ctx context.Context, p *product_domain.Product) error {
	if p.Category == nil {
		return category_domain.ErrCategoryNotFound
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
			return category_domain.ErrCategoryNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return product_domain.ErrProductAlreadyExists
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

func (r ProductRepository) DeleteProduct(ctx context.Context, productID uuid.UUID) error {

	query := "DELETE FROM products WHERE uuid = $1"

	cmdTag, err := r.dbPool.Exec(ctx, query, productID)

	if err != nil {
		return fmt.Errorf("delete product: %w", err)

	}

	if cmdTag.RowsAffected() == 0 {
		return product_domain.ErrProductNotFound
	}
	return nil

}

// UpdateProduct updates a product in the database based on its UUID.
func (r *ProductRepository) UpdateProduct(ctx context.Context, p *product_domain.Product) error {
	if p.Category == nil || p.Category.UUID == uuid.Nil {
		return category_domain.ErrCategoryNotFound
	}

	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // откатит если Commit не был

	queryProduct := `
	UPDATE products
		SET
			name = $1,
			slug = $2,
			description = $3,
			price = $4,
			stock = $5,
			category_id = (SELECT id FROM categories WHERE uuid = $6), -- если uuid не найден -> NULL
			updated_at = NOW()
		WHERE uuid = $7
	RETURNING id, uuid, name, slug, description, price, stock, category_id, is_active, created_at, updated_at
	`

	err = tx.QueryRow(ctx, queryProduct,
		p.Name, p.Slug, p.Description, p.Price, p.Stock, p.Category.UUID, p.UUID,
	).Scan(
		&p.ID, &p.UUID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock,
		&p.CategoryID, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return product_domain.ErrProductNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation - дубль slug
				return product_domain.ErrProductAlreadyExists
			case "23503": // foreign_key_violation - сломался другой FK
				return category_domain.ErrCategoryNotFound
			case "23502": // not_null_violation - SELECT вернул NULL т.к категории нет
				return category_domain.ErrCategoryNotFound
			}
		}

		return fmt.Errorf("update product: %w", err)
	}

	// обновление картинок
	if len(p.Images) > 0 {
		if _, err := tx.Exec(ctx, `DELETE FROM product_images WHERE product_id = $1`, p.ID); err != nil {
			return fmt.Errorf("delete old images: %w", err)
		}

		queryImage := `
			INSERT INTO product_images (product_id, image_url, is_main)
			VALUES ($1, $2, $3)
			RETURNING id, uuid, image_url, is_main
		`
		for i := range p.Images {
			img := &p.Images[i]
			img.ProductID = p.ID

			if err := tx.QueryRow(ctx, queryImage, img.ProductID, img.ImageURL, img.IsMain).
				Scan(&img.ID, &img.UUID, &img.ImageURL, &img.IsMain); err != nil {
				return fmt.Errorf("insert image: %w", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil { // <- ВОТ ЭТОГО НЕ ХВАТАЛО
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// AddProductImage - добавляет одну картинку к существующему товару.
// productUUID резолвится в внутренний id прямо в SQL, как и в CreateProduct.
func (r *ProductRepository) AddProductImage(ctx context.Context, productUUID uuid.UUID, imageURL string, isMain bool) (*product_domain.ProductImage, error) {
	query := `
		INSERT INTO product_images (product_id, image_url, is_main)
		SELECT p.id, $1, $2
		FROM products p
		WHERE p.uuid = $3
		RETURNING id, uuid, image_url, is_main
	`

	var img product_domain.ProductImage
	err := r.dbPool.QueryRow(ctx, query, imageURL, isMain, productUUID).
		Scan(&img.ID, &img.UUID, &img.ImageURL, &img.IsMain)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product_domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("add product image: %w", err)
	}

	return &img, nil
}

// DeleteProductImage - delete image with his ID
func (r *ProductRepository) DeleteProductImage(ctx context.Context, productUUID, imageUUID uuid.UUID) error {
	query := `
		DELETE FROM product_images
		WHERE uuid = $1
		AND product_id = (SELECT id FROM products WHERE uuid = $2)
	`

	cmdTag, err := r.dbPool.Exec(ctx, query, imageUUID, productUUID)
	if err != nil {
		return fmt.Errorf("delete product image: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return product_domain.ErrImageNotFound
	}
	return nil
}

// CountProductImages - for checking before adding new image
func (r *ProductRepository) CountProductImages(ctx context.Context, productUUID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) FROM product_images pi
		JOIN products p ON pi.product_id = p.id
		WHERE p.uuid = $1
	`
	var count int
	err := r.dbPool.QueryRow(ctx, query, productUUID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count product images: %w", err)
	}
	return count, nil
}
