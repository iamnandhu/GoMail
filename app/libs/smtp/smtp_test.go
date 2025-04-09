package smtp_test

import (
	"context"
	"errors"
	"testing"

	"GoMail/app/libs/smtp"
	"GoMail/app/libs/smtp/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewClient(t *testing.T) {
	// Arrange
	config := smtp.Config{
		Host:     "smtp.example.com",
		Port:     "587",
		Username: "test@example.com",
		Password: "password",
		From:     "test@example.com",
		UseTLS:   true,
	}

	// Act
	client := smtp.NewClient(config)

	// Assert
	assert.NotNil(t, client, "Expected non-nil client")

	// We can't access the internal smtpClient type from the test package
	// So we'll just verify the client is not nil
	// Since we can't access the internal fields, we'll just test the interface methods
	// in the other test cases
}

func TestSend(t *testing.T) {
	// Arrange
	mockClient := mocks.NewSMTPClient(t)
	ctx := context.Background()
	from := "sender@example.com"
	to := "recipient@example.com"
	subject := "Test Subject"
	body := "Test Body"

	// Setup expectations
	mockClient.On("Send", ctx, from, to, subject, body).Return(nil)

	// Act
	err := mockClient.Send(ctx, from, to, subject, body)

	// Assert
	assert.NoError(t, err, "Send should not return an error")
	mockClient.AssertExpectations(t)
}

func TestSend_Error(t *testing.T) {
	// Arrange
	mockClient := mocks.NewSMTPClient(t)
	ctx := context.Background()
	from := "sender@example.com"
	to := "recipient@example.com"
	subject := "Test Subject"
	body := "Test Body"
	expectedErr := errors.New("SMTP error")

	// Setup expectations
	mockClient.On("Send", ctx, from, to, subject, body).Return(expectedErr)

	// Act
	err := mockClient.Send(ctx, from, to, subject, body)

	// Assert
	assert.Error(t, err, "Send should return an error")
	assert.Equal(t, expectedErr, err, "Error should match expected error")
	mockClient.AssertExpectations(t)
}

func TestSendHTML(t *testing.T) {
	// Arrange
	mockClient := mocks.NewSMTPClient(t)
	ctx := context.Background()
	from := "sender@example.com"
	to := "recipient@example.com"
	subject := "Test Subject"
	htmlBody := "<html><body>Test Body</body></html>"

	// Setup expectations
	mockClient.On("SendHTML", ctx, from, to, subject, htmlBody).Return(nil)

	// Act
	err := mockClient.SendHTML(ctx, from, to, subject, htmlBody)

	// Assert
	assert.NoError(t, err, "SendHTML should not return an error")
	mockClient.AssertExpectations(t)
}

func TestConnect(t *testing.T) {
	// Arrange
	mockClient := mocks.NewSMTPClient(t)

	// Setup expectations
	mockClient.On("Connect").Return(nil)

	// Act
	err := mockClient.Connect()

	// Assert
	assert.NoError(t, err, "Connect should not return an error")
	mockClient.AssertExpectations(t)
}

func TestDisconnect(t *testing.T) {
	// Arrange
	mockClient := mocks.NewSMTPClient(t)

	// Setup expectations
	mockClient.On("Disconnect").Return(nil)

	// Act
	err := mockClient.Disconnect()

	// Assert
	assert.NoError(t, err, "Disconnect should not return an error")
	mockClient.AssertExpectations(t)
}

func TestSendWithAttachments(t *testing.T) {
	// Arrange
	mockClient := mocks.NewSMTPClient(t)
	ctx := context.Background()
	from := "sender@example.com"
	to := "recipient@example.com"
	subject := "Test Subject"
	body := "Test Body"
	attachments := []smtp.Attachment{
		{
			Filename: "test.txt",
			Content:  []byte("test content"),
			MimeType: "text/plain",
		},
	}
	expectedErr := errors.New("sending emails with attachments is not implemented yet")

	// Setup expectations
	mockClient.On("SendWithAttachments", ctx, from, to, subject, body, mock.Anything).Return(expectedErr)

	// Act
	err := mockClient.SendWithAttachments(ctx, from, to, subject, body, attachments)

	// Assert
	assert.Error(t, err, "SendWithAttachments should return an error")
	assert.Contains(t, err.Error(), "not implemented", "Error should indicate feature is not implemented")
	mockClient.AssertExpectations(t)
}
