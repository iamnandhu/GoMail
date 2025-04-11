package email

import (
	"context"
	"time"

	"GoMail/app/repository/models"
)

// SendWithAttachments sends an email with attachments
func (s *emailService) SendWithAttachments(ctx context.Context, req SendWithAttachmentsRequest) (*SendEmailResponse, error) {
	// Create a request to the SMTP client
	err := s.client.SendWithAttachments(ctx, req.From, req.To, req.Subject, req.Body, req.Attachments)
	
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
		ContentType: "multipart/mixed",
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