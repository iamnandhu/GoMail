package email

import (
	"GoMail/app/config"
	"GoMail/app/logic/email"

	"github.com/gin-gonic/gin"
)

// AddPublicRoute adds email routes that don't require authentication
func AddPublicRoute(router *gin.RouterGroup, path string) {
	// Initialize the email service directly here
	cfg := config.Get()
	
	// Create the email service
	emailService := email.NewEmailService(cfg, nil)
	handler := NewHandler(emailService)

	emailGroup := router.Group(path)
	{
		emailGroup.GET("/status", handler.getEmailStatus)
	}
}

// AddProtectedRoute adds email routes that require authentication
func AddProtectedRoute(router *gin.RouterGroup, path string) {
	// Initialize the email service directly here
	cfg := config.Get()
	
	// Create the email service
	emailService := email.NewEmailService(cfg, nil)
	handler := NewHandler(emailService)

	emailGroup := router.Group(path)
	{
		emailGroup.POST("/send", handler.sendEmail)
		emailGroup.POST("/send-html", handler.sendHTMLEmail)
		emailGroup.POST("/send-with-attachments", handler.sendEmailWithAttachments)
		emailGroup.POST("/send-bulk", handler.sendBulkEmails)
	}
} 