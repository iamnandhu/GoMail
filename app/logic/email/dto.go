package email

import (
	libSmtp "GoMail/app/libs/smtp"
)

// SendEmailRequest represents a request to send an email
type SendEmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// SendEmailResponse represents a response from sending an email
type SendEmailResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

// SendWithAttachmentsRequest represents a request to send an email with attachments
type SendWithAttachmentsRequest struct {
	From        string               `json:"from"`
	To          string               `json:"to"`
	Subject     string               `json:"subject"`
	Body        string               `json:"body"`
	Attachments []libSmtp.Attachment `json:"attachments"`
}

// SendBulkEmailRequest represents a request to send multiple emails
type SendBulkEmailRequest struct {
	Emails []BulkEmail `json:"emails"`
}

// BulkEmail represents a single email in a bulk send request
type BulkEmail struct {
	From        string               `json:"from"`
	To          string               `json:"to"`
	Subject     string               `json:"subject"`
	Body        string               `json:"body"`
	IsHTML      bool                 `json:"isHtml"`
	Attachments []libSmtp.Attachment `json:"attachments,omitempty"`
}

// SendBulkEmailResponse represents a response from sending multiple emails
type SendBulkEmailResponse struct {
	Results []EmailResult `json:"results"`
}

// EmailResult represents the result of sending a single email
type EmailResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
} 