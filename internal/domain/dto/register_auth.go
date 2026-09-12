package dto

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrEmailRequired    error = errors.New("email is required")
	ErrInvalidEmail     error = errors.New("invalid email")
	ErrPasswordRequired error = errors.New("password is required")
)

/*	 Registration		*/
type RegisterDTO struct {
	Email    string `json:"email,required"`
	Password string `json:"password,required"`
}

func (r *RegisterDTO) Validate() error {

	//  clean email from different tabs anв spaces
	r.Email = strings.TrimSpace(r.Email)

	if r.Email == "" {
		return ErrEmailRequired
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return ErrInvalidEmail
	}

	if r.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

/*	Authentication	*/

type LoginDTO struct {
	Email    string `json:"email,required"`
	Password string `json:"password,required"`
}

func (l *LoginDTO) Validate() error {
	l.Email = strings.TrimSpace(l.Email)
	if l.Password == "" {
		return ErrPasswordRequired
	}
	if _, err := mail.ParseAddress(l.Email); err != nil {
		return ErrInvalidEmail
	}

	if l.Password == "" {
		return ErrPasswordRequired
	}

	return nil
}

// AuthResponseDTO - request to user after successfully authentication
type AuthResponseDTO struct {
	AccessToken string          `json:"access_token,omitempty"`
	TokenType   string          `json:"token_type"`
	ExpiresIn   int             `json:"expires_in"`
	User        UserResponseDTO `json:"user"`
}

type UserResponseDTO struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
