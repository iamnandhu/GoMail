package email

import (
	"github.com/gin-gonic/gin"
)

// AddPublicRoute adds email routes that don't require authentication
func AddPublicRoute(router *gin.RouterGroup, path string, handler *Handler) {
	emailGroup := router.Group(path)
	{
		emailGroup.GET("/status", handler.getEmailStatus)
	}
}

// AddProtectedRoute adds email routes that require authentication
func AddProtectedRoute(router *gin.RouterGroup, path string, handler *Handler) {
	emailGroup := router.Group(path)
	{
		emailGroup.POST("/send", handler.sendEmail)
		emailGroup.POST("/send-html", handler.sendHTMLEmail)
		emailGroup.POST("/send-with-attachments", handler.sendEmailWithAttachments)
		emailGroup.POST("/send-bulk", handler.sendBulkEmails)
	}
} 