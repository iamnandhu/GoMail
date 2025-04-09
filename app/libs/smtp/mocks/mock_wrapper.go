package mocks

import (
	"context"

	"GoMail/app/libs/smtp"

	mock "github.com/stretchr/testify/mock"
)

// SMTPClientWrapper wraps the mock SMTPClient to implement all required methods
type SMTPClientWrapper struct {
	*SMTPClient
}

// NewSMTPClientWrapper creates a new wrapper around a mock SMTPClient
func NewSMTPClientWrapper(t interface {
	mock.TestingT
	Cleanup(func())
}) *SMTPClientWrapper {
	return &SMTPClientWrapper{
		SMTPClient: NewSMTPClient(t),
	}
}

// IsConnected implements the IsConnected method for the SMTPClient interface
func (w *SMTPClientWrapper) IsConnected() bool {
	// Always return true for testing
	return true
}

// Send delegates to the wrapped mock
func (w *SMTPClientWrapper) Send(ctx context.Context, from string, to string, subject string, body string) error {
	return w.SMTPClient.Send(ctx, from, to, subject, body)
}

// SendHTML delegates to the wrapped mock
func (w *SMTPClientWrapper) SendHTML(ctx context.Context, from string, to string, subject string, htmlBody string) error {
	return w.SMTPClient.SendHTML(ctx, from, to, subject, htmlBody)
}

// SendWithAttachments delegates to the wrapped mock
func (w *SMTPClientWrapper) SendWithAttachments(ctx context.Context, from string, to string, subject string, body string, attachments []smtp.Attachment) error {
	return w.SMTPClient.SendWithAttachments(ctx, from, to, subject, body, attachments)
}

// Connect delegates to the wrapped mock
func (w *SMTPClientWrapper) Connect() error {
	return w.SMTPClient.Connect()
}

// Disconnect delegates to the wrapped mock
func (w *SMTPClientWrapper) Disconnect() error {
	return w.SMTPClient.Disconnect()
} 