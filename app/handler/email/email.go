package email

import (
	"context"
	"net/http"

	"GoMail/app/logic/email"

	"github.com/gin-gonic/gin"
)

// Handler handles email-related HTTP requests
type Handler struct {
	emailService email.Email
}

// NewHandler creates a new email handler
func NewHandler(emailService email.Email) *Handler {
	return &Handler{
		emailService: emailService,
	}
}

// sendEmail handles sending a plain text email
func (h *Handler) sendEmail(c *gin.Context) {
	var req email.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.emailService.Send(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// sendHTMLEmail handles sending an HTML email
func (h *Handler) sendHTMLEmail(c *gin.Context) {
	var req email.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.emailService.SendHTML(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// sendEmailWithAttachments handles sending an email with attachments
func (h *Handler) sendEmailWithAttachments(c *gin.Context) {
	var req email.SendWithAttachmentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.emailService.SendWithAttachments(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// sendBulkEmails handles sending multiple emails
func (h *Handler) sendBulkEmails(c *gin.Context) {
	var req email.SendBulkEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.emailService.SendBulk(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// getEmailStatus returns the status of the email service
func (h *Handler) getEmailStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "operational",
	})
} 