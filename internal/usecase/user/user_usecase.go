package user_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/shodruzhoshimzoda/tojtech/internal/domain/dto"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *domain_user.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (domain_user.User, error)
	GetUserByUUID(ctx context.Context, uuid uuid.UUID) (domain_user.User, error)
}
type AuthUsercase struct {
	repo      UserRepo
	jwtSecret []byte
	ttl       time.Duration
}

func NewAuthUsercase(repo UserRepo, jwtSecret []byte, ttl time.Duration) *AuthUsercase {
	return &AuthUsercase{
		repo:      repo,
		jwtSecret: jwtSecret,
		ttl:       ttl,
	}
}

// RegisterUser - this method will register user to system
func (u *AuthUsercase) RegisterUser(ctx context.Context, userDTO dto.RegisterDTO) (dto.AuthResponseDTO, error) {
	if err := userDTO.Validate(); err != nil {
		return dto.AuthResponseDTO{}, err
	}

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}

	user := domain_user.User{
		Email:        userDTO.Email,
		PasswordHash: string(hashedPassword),
		Role:         "customer",
	}

	userUUID, err := u.repo.CreateUser(ctx, &user)

	// set user UUID

	user.UUID = userUUID

	resp := dto.AuthResponseDTO{
		TokenType: "Bearer",
		ExpiresIn: int(u.ttl.Hours()),
		User: dto.UserResponseDTO{
			UUID:  user.UUID.String(),
			Email: user.Email,
			Role:  user.Role,
		}}

	if err != nil {
		return dto.AuthResponseDTO{}, err
	}

	return resp, nil
}

// LoginUser
func (r *AuthUsercase) LoginUser(ctx context.Context, userDTO dto.LoginDTO) (dto.AuthResponseDTO, error) {
	if err := userDTO.Validate(); err != nil {
		return dto.AuthResponseDTO{}, err
	}

	user, err := r.repo.GetUserByEmail(ctx, userDTO.Email)
	if err != nil {
		return dto.AuthResponseDTO{}, domain_user.ErrUserNotFound
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDTO.Password)); err != nil {
		return dto.AuthResponseDTO{}, domain_user.ErrInvalidEmailOrPassword
	}

	// JWT token generating

	claims := jwt.MapClaims{
		"sub":  user.UUID.String(),
		"exp":  time.Now().Add(r.ttl).Unix(),
		"role": user.Role,
		"iat":  time.Now().Unix(),
	}

	// sign token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(r.jwtSecret)
	if err != nil {

		return dto.AuthResponseDTO{}, fmt.Errorf("failed to sign token: %w", err)
	}

	resp := dto.AuthResponseDTO{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int(r.ttl.Hours()),
		User: dto.UserResponseDTO{
			UUID:  user.UUID.String(),
			Email: user.Email,
			Role:  user.Role,
		}}

	return resp, nil
}
