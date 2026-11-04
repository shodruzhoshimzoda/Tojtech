package dto

import (
	"errors"
	"net/mail"
	"strings"

	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
)

var (
	ErrEmailRequired    error = errors.New("email is required")
	ErrPasswordRequired error = errors.New("password is required")
)

/*	 Registration		*/
type RegisterDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *RegisterDTO) Validate() error {

	//  clean email from different tabs anв spaces
	r.Email = strings.TrimSpace(r.Email)

	if r.Email == "" {
		return ErrEmailRequired
	}

	if _, err := mail.ParseAddress(r.Email); err != nil {
		return domain_user.ErrInvalidEmailOrPassword
	}

	if r.Password == "" {
		return ErrPasswordRequired
	}

	// strict rules for a NEW password: length, allowed characters, letter case, digit
	if err := domain_user.ValidatePassword(r.Password); err != nil {
		
    return err
}

	return nil
}

/*	Authentication	*/

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginDTO) Validate() error {

	l.Email = strings.TrimSpace(l.Email)
	if l.Password == "" {
		return ErrPasswordRequired
	}
	if l.Email == "" {
		return ErrEmailRequired
	}

	if _, err := mail.ParseAddress(l.Email); err != nil {
		return domain_user.ErrInvalidEmailOrPassword
	}

	return nil
}

// AuthResponseDTO - request to user after successfully authentication
type AuthResponseDTO struct {
	AccessToken string          `json:"access_token,omitempty"`
	RefreshToken string          `json:"refresh_token,omitempty"`
	TokenType   string          `json:"token_type"`
	ExpiresIn   int             `json:"expires_in"`
	User        UserResponseDTO `json:"user"`
}
type RefreshRequestDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type UserResponseDTO struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
