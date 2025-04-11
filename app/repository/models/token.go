package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RevokedToken represents a revoked JWT token in the database
type RevokedToken struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TokenID     string             `bson:"token_id" json:"token_id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	RevokedAt   time.Time          `bson:"revoked_at" json:"revoked_at"`
	ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
	Reason      string             `bson:"reason,omitempty" json:"reason,omitempty"`
}