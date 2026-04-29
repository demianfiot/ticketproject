package repository

import (
	"context"

	domain1 "github.com/demianfiot/ticketproject/auth-service/internal"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(ctx context.Context, user domain1.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (domain1.User, error)
	GetUserByID(ctx context.Context, userID int) (domain1.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
