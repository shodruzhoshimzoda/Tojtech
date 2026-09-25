package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	domain_user "github.com/shodruzhoshimzoda/tojtech/internal/domain/user"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}

}

// CreateUser - add user to db and register him
func (r *UserRepository) CreateUser(ctx context.Context, user *domain_user.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING uuid
	`

	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRow(ctx, query,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.UUID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return uuid.Nil, domain_user.ErrUserAlreadyExists
		}
		return uuid.Nil, err
	}

	return user.UUID, nil

}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain_user.User, error) {
	query := `SELECT uuid, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`

	var user domain_user.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain_user.User{}, domain_user.ErrUserNotFound
		}
		return domain_user.User{}, err
	}
	return user, nil

}

func (r *UserRepository) GetUserByUUID(ctx context.Context, uuid uuid.UUID) (domain_user.User, error) {

	query := `
		SELECT uuid, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE uuid = $1`
	var user domain_user.User
	err := r.db.QueryRow(ctx, query, uuid).Scan(
		&user.UUID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain_user.User{}, domain_user.ErrUserNotFound
		}
		return domain_user.User{}, err
	}
	return user, nil

}
