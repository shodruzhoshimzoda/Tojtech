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
func (u *AuthUsercase) RegisterUser(ctx context.Context, dto dto.RegisterDTO) (uuid.UUID, error) {
	if err := dto.Validate(); err != nil {
		return uuid.Nil, err
	}

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	user := &domain_user.User{
		Email:        dto.Email,
		PasswordHash: string(hashedPassword),
		Role:         "customer",
	}

	userUUid, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUid, nil
}

// LoginUser
func (r *AuthUsercase) LoginUser(ctx context.Context, dto dto.LoginDTO) (string, error) {
	if err := dto.Validate(); err != nil {
		return "", err
	}

	user, err := r.repo.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return "", domain_user.ErrInvalidEmailOrPassword
	}

	// check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.Password)); err != nil {
		return "", domain_user.ErrInvalidEmailOrPassword
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
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}
