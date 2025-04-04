package repository

import (
	"context"
)

// Repository is the interface that all repositories should implement
type Repository interface {
	// Common methods that all repositories should implement
	Connect() error
	Disconnect() error
}

// EmailRepository defines operations for email persistence
type EmailRepository interface {
	Repository
	SaveEmail(ctx context.Context, email interface{}) error
	FindEmails(ctx context.Context, filter interface{}) ([]interface{}, error)
	FindEmailByID(ctx context.Context, id string) (interface{}, error)
} 