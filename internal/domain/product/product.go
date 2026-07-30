package domain

import (
	"time"
)

type Product struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       int64      `json:"price"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
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
