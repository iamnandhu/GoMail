package email

import (
	"context"
	"fmt"
	"time"

	"GoMail/app/config"
	"GoMail/app/libs/smtp"
	"GoMail/app/repository"
	"GoMail/app/repository/models"
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
	fmt.Printf("  TLS Enabled: %v\n", cfg.SMTP.UseStartTLS)
	fmt.Printf("  Max Concurrent: %v\n", cfg.SMTP.MaxConcurrent)
	
	// Set a default for MaxConcurrent if not specified
	maxConcurrent := 10
	if cfg.SMTP.MaxConcurrent > 0 {
		maxConcurrent = cfg.SMTP.MaxConcurrent
	}
	
	// Create SMTP config from application config
	smtpConfig := smtp.Config{
		Host:               cfg.SMTP.Host,
		Port:               cfg.SMTP.Port,
		Username:           cfg.SMTP.Username,
		Password:           cfg.SMTP.Password,
		From:               cfg.SMTP.From,
		UseTLS:             !cfg.SMTP.UseStartTLS, // Use TLS if not using StartTLS
		StartTLS:           cfg.SMTP.UseStartTLS,  // Use StartTLS from config
		InsecureSkipVerify: true,                  // Skip verification for testing
		ConnectTimeout:     10 * time.Second,
		PoolSize:           5,
		RetryAttempts:      3,
		RetryDelay:         2 * time.Second,
		MaxConcurrent:      maxConcurrent,
	}
	
	// Debug: Print SMTP config after conversion
	fmt.Printf("DEBUG: Created SMTP config with Host=%s, Port=%s, UseTLS=%v, StartTLS=%v, MaxConcurrent=%v\n", 
		smtpConfig.Host, smtpConfig.Port, smtpConfig.UseTLS, smtpConfig.StartTLS, smtpConfig.MaxConcurrent)
	
	// Create SMTP client with config
	client := smtp.NewClient(smtpConfig)
	
	return &emailService{
		client: client,
		repo:   repo,
		config: cfg,
	}
}

// logEmailAttempt logs an email attempt asynchronously
func (s *emailService) logEmailAttempt(logData *models.EmailLog) {
	// Only proceed if repository is available
	if s.repo == nil {
		return
	}
	
	// Log the email asynchronously 
	go func() {
		// Create a new context for the async operation
		asyncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Log the email
		err := s.repo.SaveEmailLog(asyncCtx, logData)
		if err != nil {
			// Just log the error, don't propagate it
			fmt.Printf("ERROR: Failed to log email: %v\n", err)
		}
	}()
} 