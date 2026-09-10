package domain_product

import (
	"errors"
	"net/url"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	domain_category "github.com/shodruzhoshimzoda/tojtech/internal/domain/category"
)

var (
	ErrEmptyProductName    = errors.New("product name cannot be empty")
	ErrInvalidProductName  = errors.New("product name must be less than 100 characters")
	ErrEmptyProductSlug    = errors.New("product slug cannot be empty")
	ErrInvalidProductPrice = errors.New("product price must be greater than zero")
	ErrNegativeStock       = errors.New("product stock cannot be negative")
	ErrLongDescription     = errors.New("product description must be less than 500 characters")

	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product with this slug already exists")
)

type Product struct {
	ID          int64                     `db:"id" json:"-"`
	UUID        uuid.UUID                 `db:"uuid" json:"uuid"`
	Name        string                    `db:"name" json:"name"`
	Slug        string                    `db:"slug" json:"slug"`
	Description string                    `db:"description" json:"description,omitempty"`
	Price       float64                   `db:"price" json:"price"`
	Stock       int                       `db:"stock" json:"stock"`
	CategoryID  *int64                    `db:"category_id" json:"-"` // внутренний id, скрыли
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
	ImageURL  string    `db:"image_url" json:"image_url"`
	IsMain    bool      `db:"is_main" json:"is_main"`
}

// internal/domain/product/product.go
var (
	ErrImageNotFound   = errors.New("image not found")
	ErrEmptyImageURL   = errors.New("image url cannot be empty")
	ErrInvalidImageURL = errors.New("image url is not a valid url")
	ErrTooManyImages   = errors.New("product cannot have more than 10 images")
)

func ValidateImageURL(rawURL string) error {
	if rawURL == "" {
		return ErrEmptyImageURL
	}
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return ErrInvalidImageURL
	}
	return nil
}

// type ProductDTO struct {
// 	ID				int64 		`json:"id"`
// 	Name 			string 		`json:"name"`
// 	Description		string		`json:"description"`
// 	Price			int64		`json:"price"`

// }

// Validate - this method will check each fields for our product
func (p *Product) Validate() error {

	if p.Name == "" {
		return ErrEmptyProductName
	}

	if utf8.RuneCountInString(p.Name) > 100 {
		return ErrInvalidProductName
	}
	if p.Price < 0 {
		return ErrInvalidProductPrice
	}
	if p.Slug == "" {
		return ErrEmptyProductSlug
	}
	if p.Stock < 0 {
		return ErrNegativeStock
	}
	if len(p.Description) > 500 {
		return ErrLongDescription
	}

	return nil

}
