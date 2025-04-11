package handler

import (
	"GoMail/app/handler/auth"
	"GoMail/app/handler/email"
	"GoMail/app/middleware"

	"github.com/gin-gonic/gin"
)

// InitPublicRoutes initializes routes that don't require authentication
func InitPublicRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	
	// Register public routes
	auth.AddRoute(api, "/auth")
}

// InitProtectedRoutes initializes routes that require authentication
func InitProtectedRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	api.Use(middleware.VerifyAuthToken())
	
	// Register protected routes
	email.AddProtectedRoute(api, "/email")
}