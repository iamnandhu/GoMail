package email

import (
	"context"
	"sync"
	"time"

	"GoMail/app/repository/models"
)

// SendBulk sends multiple emails concurrently
func (s *emailService) SendBulk(ctx context.Context, req SendBulkEmailRequest) (*SendBulkEmailResponse, error) {
	// Initialize a slice to store the results
	results := make([]EmailResult, len(req.Emails))

	// Create a wait group to wait for all emails to be sent
	var wg sync.WaitGroup
	wg.Add(len(req.Emails))

	// Set a limit on the number of concurrent requests
	// Default to 5 concurrent connections
	maxConcurrent := 5
	
	// Create a semaphore using a channel
	semaphore := make(chan struct{}, maxConcurrent)
	
	// Send emails concurrently
	for i, email := range req.Emails {
		go func(idx int, email BulkEmail) {
			defer wg.Done()
			
			// Acquire a semaphore slot
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			var err error
			
			// Determine the content type
			contentType := "text/plain"
			
			// Send the email based on its type
			if email.IsHTML {
				contentType = "text/html"
				err = s.client.SendHTML(ctx, email.From, email.To, email.Subject, email.Body)
			} else if len(email.Attachments) > 0 {
				contentType = "multipart/mixed"
				err = s.client.SendWithAttachments(ctx, email.From, email.To, email.Subject, email.Body, email.Attachments)
			} else {
				err = s.client.Send(ctx, email.From, email.To, email.Subject, email.Body)
			}
			
			// Create success/error response
			success := err == nil
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			
			// Store the result
			results[idx] = EmailResult{
				Success: success,
				Error:   errMsg,
			}
			
			// Create email log
			emailLog := &models.EmailLog{
				From:        email.From,
				To:          email.To,
				Subject:     email.Subject,
				ContentType: contentType,
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
		}(i, email)
	}
	
	// Wait for all emails to be sent
	wg.Wait()
	
	// Return the results
	return &SendBulkEmailResponse{
		Results: results,
	}, nil
} 