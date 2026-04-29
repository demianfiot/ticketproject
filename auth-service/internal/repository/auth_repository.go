package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain1 "github.com/demianfiot/ticketproject/auth-service/internal"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrUserNotFound   = errors.New("user not found")
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user domain1.User) (int, error) {
	var id int

	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, ErrDuplicateEmail
		}

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *AuthPostgres) GetUserByEmail(ctx context.Context, email string) (domain1.User, error) {
	var user domain1.User

	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain1.User{}, ErrUserNotFound
		}

		return domain1.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *AuthPostgres) GetUserByID(ctx context.Context, userID int) (domain1.User, error) {
	var user domain1.User

	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain1.User{}, ErrUserNotFound
		}

		return domain1.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}
