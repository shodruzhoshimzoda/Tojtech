package domain_product

import "errors"

var (
	ErrInvalidProductName  = errors.New("The length of the product name must be less than 100 characters and greater than 0.")
	ErrEmptyProductName    = errors.New("Product name could no be empty")
	ErrInvalidProductPrice = errors.New("the product price is incorrect")
	ErrLongDescription     = errors.New("The product description is too long")

	ErrProductNotFound = errors.New("product not  found")
)
