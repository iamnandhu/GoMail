package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmailLog represents a log of an email that was sent
type EmailLog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	From        string             `bson:"from" json:"from"`
	To          string             `bson:"to" json:"to"`
	Subject     string             `bson:"subject" json:"subject"`
	ContentType string             `bson:"content_type" json:"content_type"`
	Success     bool               `bson:"success" json:"success"`
	SentAt      time.Time          `bson:"sent_at" json:"sent_at"`
	Error       string             `bson:"error,omitempty" json:"error,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
} 