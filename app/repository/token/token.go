package token

import (
	"GoMail/app/repository/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tokenCollection = "revoked_tokens"
)

var (
	ErrTokenAlreadyRevoked = errors.New("token already revoked")
)

// Repository defines the token repository interface
type Repository interface {
	// RevokeToken adds a token to the revoked tokens collection
	RevokeToken(ctx context.Context, token *models.RevokedToken) error
	
	// IsTokenRevoked checks if a token is in the revoked tokens collection
	IsTokenRevoked(ctx context.Context, tokenID string) (bool, error)
	
	// CleanupExpiredTokens removes tokens that have expired beyond the TTL
	CleanupExpiredTokens(ctx context.Context, beforeTime time.Time) (int64, error)
	
	// CreateTokenRevokedIndex creates an index for token revocation performance
	CreateTokenRevokedIndex(ctx context.Context) error
}

// TokenRepository is the MongoDB implementation of Repository
type TokenRepository struct {
	db *mongo.Database
}

// New creates a new token repository
func New(db *mongo.Database) Repository {
	return &TokenRepository{
		db: db,
	}
} 