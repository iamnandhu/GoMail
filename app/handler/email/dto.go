package email

// This file would contain any handler-specific DTOs
// Most of the DTOs are already defined in the logic/email package

// StatusResponse represents a response for checking email service status
type StatusResponse struct {
	Status string `json:"status"`
} 