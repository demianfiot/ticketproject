package service

import (
	"context"

	domain "github.com/demianfiot/ticketproject/auth-service/internal"
	"github.com/demianfiot/ticketproject/auth-service/internal/repository"
)

type Authorization interface {
	CreateUser(ctx context.Context, input CreateUserInput) (int, error)
	GenerateToken(ctx context.Context, email, password string) (string, error)
	ParseToken(ctx context.Context, accessToken string) (int, error)
	GetUserByID(ctx context.Context, userID int) (domain.User, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository, jwtSecret string, tokenTTLHours int) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization, jwtSecret, tokenTTLHours),
	}
}
