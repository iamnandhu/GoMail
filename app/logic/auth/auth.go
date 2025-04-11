package auth

import (
	"GoMail/app/config"
	"GoMail/app/repository"
	"GoMail/app/repository/models"
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidToken       = errors.New("invalid token")
)

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// Service defines the authentication service
type Service interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, user *models.User, password string) error
	VerifyToken(tokenString string) (*Claims, error)
	RevokeToken(ctx context.Context, tokenString string) error
}

type service struct {
	repo   repository.Repository
	config *config.Config
}

// New creates a new authentication service
func New(repo repository.Repository, cfg *config.Config) Service {
	return &service{
		repo:   repo,
		config: cfg,
	}
} 