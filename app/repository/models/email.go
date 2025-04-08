package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmailStatus represents the status of an email
type EmailStatus string

const (
	// EmailStatusPending indicates the email is pending to be sent
	EmailStatusPending EmailStatus = "pending"
	// EmailStatusSent indicates the email was successfully sent
	EmailStatusSent EmailStatus = "sent"
	// EmailStatusFailed indicates the email failed to send
	EmailStatusFailed EmailStatus = "failed"
)

// Email represents an email document in the database
type Email struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	From      string             `bson:"from" json:"from"`
	To        string             `bson:"to" json:"to"`
	Subject   string             `bson:"subject" json:"subject"`
	Body      string             `bson:"body" json:"body"`
	IsHTML    bool               `bson:"is_html" json:"is_html"`
	Status    EmailStatus        `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	SentAt    *time.Time         `bson:"sent_at,omitempty" json:"sent_at,omitempty"`
	Error     string             `bson:"error,omitempty" json:"error,omitempty"`
} 