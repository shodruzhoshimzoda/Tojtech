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
	ErrInvalidRole            = errors.New("role must be either admin or customer")
	ErrInvalidRefreshToken    = errors.New("invalid or expired refresh token")
	ErrCannotChangeOwnRole    = errors.New("you cannot change your own role")

	// for user profile
	ErrWrongCurrentPassword   = errors.New("current password is wrong")
	ErrSamePassword			  = errors.New("The password must be different from current")

)

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

type User struct {
	ID           int64     `json:"id"`
	UUID         uuid.UUID `json:"uuid"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}


func IsValidRole(role string) bool {
	
	return role == RoleAdmin || role==RoleCustomer
}