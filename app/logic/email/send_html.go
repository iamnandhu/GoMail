package email

import (
	"context"
	"time"

	"GoMail/app/repository/models"
)

// SendHTML sends an HTML email
func (s *emailService) SendHTML(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error) {
	// Create a request to the SMTP client
	err := s.client.SendHTML(ctx, req.From, req.To, req.Subject, req.Body)
	
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
		ContentType: "text/html",
		SentAt:      time.Now(),
		Success:     success,
		Error:       errMsg,
		CreatedAt:   time.Now(),
	}
	
	// Log the email asynchronously
	s.logEmailAttempt(emailLog)
	
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