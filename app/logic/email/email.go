package email

import (
	"context"

	"GoMail/app/config"
	"GoMail/app/libs/smtp"
	"GoMail/app/repository"
)

// Email defines the interface for email operations
type Email interface {
	// Send sends a plain text email
	Send(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error)
	
	// SendHTML sends an HTML email
	SendHTML(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error)
	
	// SendWithAttachments sends an email with attachments
	SendWithAttachments(ctx context.Context, req SendWithAttachmentsRequest) (*SendEmailResponse, error)
	
	// SendBulk sends multiple emails concurrently
	SendBulk(ctx context.Context, req SendBulkEmailRequest) (*SendBulkEmailResponse, error)
}

// emailService implements the Email interface
type emailService struct {
	client smtp.SMTPClient
	repo   repository.Repository
	config *config.Config
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config, client smtp.SMTPClient, repo repository.Repository) Email {
	return &emailService{
		client: client,
		repo:   repo,
		config: cfg,
	}
} 