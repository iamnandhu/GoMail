package logic

import (
	"context"
)

// EmailService defines operations for email handling
type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SaveEmailRecord(ctx context.Context, emailData interface{}) error
	GetEmailHistory(ctx context.Context) ([]interface{}, error)
} 