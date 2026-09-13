package user_usecase

import (
	"context"
	"errors"
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
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}

	// set user UUID

	user.UUID = userUUID

	return u.GenerateToken(&user)
}

// LoginUser
func (u *AuthUsercase) LoginUser(ctx context.Context, userDTO dto.LoginDTO) (dto.AuthResponseDTO, error) {
	if err := userDTO.Validate(); err != nil {
		return dto.AuthResponseDTO{}, err
	}

	user, err := u.repo.GetUserByEmail(ctx, userDTO.Email)
	if err != nil {
		if errors.Is(err, domain_user.ErrUserNotFound) {
			return dto.AuthResponseDTO{}, domain_user.ErrInvalidEmailOrPassword

		}

		return dto.AuthResponseDTO{}, err
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userDTO.Password)); err != nil {
		return dto.AuthResponseDTO{}, domain_user.ErrInvalidEmailOrPassword
	}

	// sign token
	return u.GenerateToken(&user)
}
func (u *AuthUsercase) GenerateToken(user *domain_user.User) (dto.AuthResponseDTO, error) {
	claims := jwt.MapClaims{
		"sub":  user.UUID.String(),
		"role": user.Role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(u.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return dto.AuthResponseDTO{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return dto.AuthResponseDTO{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   int(u.ttl.Seconds()), // Возвращаем время жизни в секундах
		User: dto.UserResponseDTO{
			UUID:  user.UUID.String(),
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
