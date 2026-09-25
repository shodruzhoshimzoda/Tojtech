package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidEmailOrPassword = errors.New("invalid user or password")
	ErrInvalidRole = errors.New("role must be either admin or customer")
)

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

type User struct {
	ID           int64     `json:"id"`
	UUID         uuid.UUID `json:"uuid"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
