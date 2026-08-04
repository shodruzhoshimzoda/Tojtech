package domain_category

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          int64     `db:"id" json:"-"` // скрыли из json
	UUID        uuid.UUID `db:"uuid" json:"uuid"`
	Name        string    `db:"name" json:"name"`
	Slug        string    `db:"slug" json:"slug"`
	Description string    `db:"description" json:"description,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
