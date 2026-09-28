package user_usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
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
	StoreRefreshToken(ctx context.Context, userUUID uuid.UUID, tokenHash string, expiresAt time.Time) error
	GetActiveRefreshTokenUserUUID(ctx context.Context, tokenHash string) (uuid.UUID, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
}

type AuthUsecase struct {
	repo       UserRepo
	jwtSecret  []byte
	ttl        time.Duration
	refreshTTL time.Duration
}

func NewAuthUsecase(repo UserRepo, jwtSecret []byte, ttl, refreshTTL time.Duration) *AuthUsecase {
	return &AuthUsecase{
		repo:       repo,
		jwtSecret:  jwtSecret,
		ttl:        ttl,
		refreshTTL: refreshTTL,
	}
}

// RegisterUser - this method will register user to system
func (u *AuthUsecase) RegisterUser(ctx context.Context, userDTO dto.RegisterDTO) (dto.AuthResponseDTO, error) {
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

	return u.issueTokens(ctx, &user)
}

// LoginUser
func (u *AuthUsecase) LoginUser(ctx context.Context, userDTO dto.LoginDTO) (dto.AuthResponseDTO, error) {
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
	return u.issueTokens(ctx, &user)
}

func (u *AuthUsecase) RefreshToken(ctx context.Context, rawRefreshToken string) (dto.AuthResponseDTO, error) {
	tokenHash := hashRefreshToken(rawRefreshToken)

	userUUID, err := u.repo.GetActiveRefreshTokenUserUUID(ctx, tokenHash)
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}

	if err := u.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return dto.AuthResponseDTO{}, domain_user.ErrInvalidRefreshToken
	}

	user, err := u.repo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, domain_user.ErrUserNotFound) {
			return dto.AuthResponseDTO{}, domain_user.ErrInvalidRefreshToken
		}
		return dto.AuthResponseDTO{}, err
	}

	return u.issueTokens(ctx, &user)
}


func (u *AuthUsecase) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := hashRefreshToken(rawRefreshToken)
	_ = u.repo.RevokeRefreshToken(ctx, tokenHash)
	return nil
}

func (u *AuthUsecase) issueTokens(ctx context.Context, user *domain_user.User) (dto.AuthResponseDTO, error) {
	accessToken, err := u.signAccessToken(user)
	if err != nil {
		return dto.AuthResponseDTO{}, err
	}

	rawRefreshToken, tokenHash, err := generateRefreshToken()
	if err != nil {
		return dto.AuthResponseDTO{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(u.refreshTTL)
	if err := u.repo.StoreRefreshToken(ctx, user.UUID, tokenHash, expiresAt); err != nil {
		return dto.AuthResponseDTO{}, err
	}

	return dto.AuthResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken, // сырой токен уходит клиенту, хеш остаётся только в БД
		TokenType:    "Bearer",
		ExpiresIn:    int(u.ttl.Seconds()),
		User: dto.UserResponseDTO{
			UUID:  user.UUID.String(),
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (u *AuthUsecase) signAccessToken(user *domain_user.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  user.UUID.String(),
		"role": user.Role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(u.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil
}

func generateRefreshToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw = hex.EncodeToString(b)
	return raw, hashRefreshToken(raw), nil
}

func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}