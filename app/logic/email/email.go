package email

import (
	"context"
	"fmt"
	"time"

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
func NewEmailService(cfg *config.Config, repo repository.Repository) Email {
	// Debug: Print SMTP config from config object
	fmt.Printf("DEBUG: Creating email service with SMTP config:\n")
	fmt.Printf("  Host: %s\n", cfg.SMTP.Host)
	fmt.Printf("  Port: %s\n", cfg.SMTP.Port)
	fmt.Printf("  Username: %s\n", cfg.SMTP.Username)
	fmt.Printf("  From: %s\n", cfg.SMTP.From)
	fmt.Printf("  TLS Enabled: %v\n", cfg.SMTP.TLSEnable)
	
	// Create SMTP config from application config
	smtpConfig := smtp.Config{
		Host:               cfg.SMTP.Host,
		Port:               cfg.SMTP.Port,
		Username:           cfg.SMTP.Username,
		Password:           cfg.SMTP.Password,
		From:               cfg.SMTP.From,
		UseTLS:             false,     // Don't use immediate TLS
		StartTLS:           true,      // Use StartTLS for encryption
		InsecureSkipVerify: true,      // Skip verification for testing
		ConnectTimeout:     10 * time.Second,
		PoolSize:           5,
		RetryAttempts:      3,
		RetryDelay:         2 * time.Second,
		MaxConcurrent:      10,
	}
	
	// Debug: Print SMTP config after conversion
	fmt.Printf("DEBUG: Created SMTP config with Host=%s, Port=%s, UseTLS=%v, StartTLS=%v\n", 
		smtpConfig.Host, smtpConfig.Port, smtpConfig.UseTLS, smtpConfig.StartTLS)
	
	// Create SMTP client with config
	client := smtp.NewClient(smtpConfig)
	
	return &emailService{
		client: client,
		repo:   repo,
		config: cfg,
	}
} 