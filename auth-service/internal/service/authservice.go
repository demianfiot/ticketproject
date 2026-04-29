package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	domain "github.com/demianfiot/ticketproject/auth-service/internal"
	"github.com/demianfiot/ticketproject/auth-service/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

type tokenClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo      repository.Authorization
	jwtSecret string
	tokenTTL  time.Duration
}

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

func NewAuthService(repo repository.Authorization, jwtSecret string, tokenTTLHours int) *AuthService {
	if tokenTTLHours <= 0 {
		tokenTTLHours = 12
	}

	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		tokenTTL:  time.Duration(tokenTTLHours) * time.Hour,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input CreateUserInput) (int, error) {
	passwordHash, err := s.generatePasswordHash(input.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	user := domain.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: passwordHash,
	}

	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return 0, ErrUserAlreadyExists
		}

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}

		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ParseToken(ctx context.Context, accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	return claims.UserID, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID int) (domain.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *AuthService) generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
