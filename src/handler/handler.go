package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(router *gin.RouterGroup) {
	// Register email handlers
	email := router.Group("/email")
	{
		email.POST("/send", sendEmail)
	}
}

// sendEmail handles the email sending request
func sendEmail(c *gin.Context) {
	// This is a placeholder for the actual implementation
	c.JSON(200, gin.H{
		"message": "Email sent successfully",
	})
} 