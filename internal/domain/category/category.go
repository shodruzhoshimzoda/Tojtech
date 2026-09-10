package category_domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidCategoryName = errors.New("the length of the product name must be less than 100 characters and greater than 0")
	ErrEmptyCategoryName   = errors.New("category name could no be empty")
	ErrSlugEmpty           = errors.New("slug could not be emty")

	ErrLongDescription       = errors.New("the product description is too long")
	ErrInvalidUUID           = errors.New("invalud uuid")
	ErrCategoryAlreadyExists = errors.New("the category already exists")
	ErrCategoryNotFound      = errors.New("category not  found")
)

type Category struct {
	ID          int64     `db:"id" json:"-"` // скрыли из json
	UUID        uuid.UUID `db:"uuid" json:"uuid"`
	Name        string    `db:"name" json:"name"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

func (p *Category) Validate() error {

	if p.Name == "" {
		return ErrEmptyCategoryName
	}

	if len(p.Name) > 100 {
		return ErrInvalidCategoryName
	}

	if len(p.Description) > 500 {
		return ErrLongDescription
	}

	return nil

}
