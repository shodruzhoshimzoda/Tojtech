package domain_product

import (
	"time"

	"github.com/google/uuid"
	domain_category "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
)

type Product struct {
	ID          int64                     `db:"id" json:"-"`
	UUID        uuid.UUID                 `db:"uuid" json:"uuid"`
	Name        string                    `db:"name" json:"name"`
	Slug        string                    `db:"slug" json:"slug"`
	Description string                    `db:"description" json:"description,omitempty"`
	Price       float64                   `db:"price" json:"price"`
	Stock       int                       `db:"stock" json:"stock"`
	CategoryID  int64                     `db:"category_id" json:"-"` // внутренний id, скрыли
	Category    *domain_category.Category `json:"category,omitempty"` // вложим категорию в ответ
	IsActive    bool                      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time                 `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time                 `db:"updated_at" json:"updated_at"`
	Images      []ProductImage            `json:"images,omitempty"`
}

type ProductImage struct {
	ID        int64     `db:"id" json:"-"`
	UUID      uuid.UUID `db:"uuid" json:"uuid"`
	ProductID int64     `db:"product_id" json:"-"`
	ImageURL  string    `db:"image_url" json:"url"`
	IsMain    bool      `db:"is_main" json:"is_main"`
}

// type ProductDTO struct {
// 	ID				int64 		`json:"id"`
// 	Name 			string 		`json:"name"`
// 	Description		string		`json:"description"`
// 	Price			int64		`json:"price"`

// }

// validate - this method will check each fields for our product
func (p *Product) Validate() error {

	if p.Name == "" {
		return ErrEmptyProductName
	}

	if len(p.Name) > 100 {
		return ErrInvalidProductName
	}

	if p.Price < 0 {
		return ErrInvalidProductPrice
	}

	if len(p.Description) > 500 {
		return ErrLongDescription
	}

	return nil

}
