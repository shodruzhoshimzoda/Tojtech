package domain

import (
	"errors"
	"time"
)



// Errors for products
var (
	ErrInvalidProductName = errors.New("The length of the product name must be less than 100 characters and greater than 0.")
	ErrEmptyProductName = errors.New("Product name could no be empty")
	ErrInvalidProductPrice = errors.New("the product price is incorrect")
	ErrLongDescription = errors.New("The product description is too long")
)



type Product struct {
	ID				int64 		`json:"id"`
	Name 			string 		`json:"name"`
	Description		string		`json:"description"`
	Price			int64		`json:"price"`
	CreatedAt 		*time.Time	`json:"created_at,omitempty"`
	UpdatedAt 		*time.Time	`json:"updated_at,omitempty"`
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

	if  len(p.Name) > 100 {
		return  ErrInvalidProductName
	}

	if p.Price < 0 {
		return ErrInvalidProductPrice
	}

	if len(p.Description) > 500 {
		return ErrLongDescription
	}

	return nil

}


