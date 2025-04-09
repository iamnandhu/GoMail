package email

import (
	"context"
	"time"

	"GoMail/app/repository/models"
)

// Send sends a plain text email
func (s *emailService) Send(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error) {
	// Create a request to the SMTP client
	err := s.client.Send(ctx, req.From, req.To, req.Subject, req.Body)
	
	// Create success/error response
	success := err == nil
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	
	// Create email log
	emailLog := &models.EmailLog{
		From:        req.From,
		To:          req.To,
		Subject:     req.Subject,
		ContentType: "text/plain",
		SentAt:      time.Now(),
		Success:     success,
		Error:       errMsg,
		CreatedAt:   time.Now(),
	}
	
	// Log the email asynchronously
	go func() {
		// Create a new context for the async operation
		asyncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		// Log the email
		_ = s.repo.SaveEmailLog(asyncCtx, emailLog)
	}()
	
	// Return the response
	if err != nil {
		return &SendEmailResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}
	
	return &SendEmailResponse{
		Success: true,
	}, nil
} 