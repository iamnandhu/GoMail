package token

import (
	"GoMail/app/repository/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// RevokeToken adds a token to the revoked tokens collection
func (r *TokenRepository) RevokeToken(ctx context.Context, token *models.RevokedToken) error {
	collection := r.db.Collection(tokenCollection)
	
	// Set default values if not provided
	if token.RevokedAt.IsZero() {
		token.RevokedAt = time.Now()
	}
	
	if token.ID.IsZero() {
		token.ID = primitive.NewObjectID()
	}
	
	// Check if token is already revoked
	exists, err := r.IsTokenRevoked(ctx, token.TokenID)
	if err != nil {
		return err
	}
	
	if exists {
		return ErrTokenAlreadyRevoked
	}
	
	// Insert the token
	_, err = collection.InsertOne(ctx, token)
	return err
}

// IsTokenRevoked checks if a token is in the revoked tokens collection
func (r *TokenRepository) IsTokenRevoked(ctx context.Context, tokenID string) (bool, error) {
	collection := r.db.Collection(tokenCollection)
	
	filter := bson.M{"token_id": tokenID}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

// CleanupExpiredTokens removes tokens that have expired beyond the TTL
func (r *TokenRepository) CleanupExpiredTokens(ctx context.Context, beforeTime time.Time) (int64, error) {
	collection := r.db.Collection(tokenCollection)
	
	filter := bson.M{"expires_at": bson.M{"$lt": beforeTime}}
	result, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	
	return result.DeletedCount, nil
}

// CreateTokenRevokedIndex creates an index for token revocation performance
func (r *TokenRepository) CreateTokenRevokedIndex(ctx context.Context) error {
	collection := r.db.Collection(tokenCollection)
	
	// Create indexes
	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "token_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "expires_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	}
	
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	return err
} 