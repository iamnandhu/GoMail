package smtp

import "time"

// Attachment represents an email attachment
type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// Config holds configuration for the SMTP client
type Config struct {
	Host               string
	Port               string
	Username           string
	Password           string
	From               string
	UseTLS             bool
	StartTLS           bool
	InsecureSkipVerify bool
	ConnectTimeout     time.Duration
	PoolSize           int
	RetryAttempts      int
	RetryDelay         time.Duration
	MaxConcurrent      int
}

// EmailRequest represents a request to send an email
type EmailRequest struct {
	From        string
	To          string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []Attachment
}

// EmailResponse represents a response from sending an email
type EmailResponse struct {
	Success bool
	Error   string
}
