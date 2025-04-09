package mocks

import (
	"GoMail/app/libs/smtp"
	"context"
)

// IsConnected provides a mock function for checking connection status
func (m *SMTPClient) IsConnected() bool {
	// For testing we can just return true by default
	return true
}

// SendBulk provides a mock function for the bulk send method
func (m *SMTPClient) SendBulk(ctx context.Context, requests []smtp.EmailRequest) []smtp.EmailResponse {
	// This method is not used directly in our tests since we mock Send and SendHTML separately
	return []smtp.EmailResponse{}
} 