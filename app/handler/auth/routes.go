package auth

import (
	"GoMail/app/config"
	"GoMail/app/logic/auth"
	"GoMail/app/middleware"
	"GoMail/app/repository"

	"github.com/gin-gonic/gin"
)

// AddRoute adds auth routes to the given router
func AddRoute(router *gin.RouterGroup, path string, repo repository.Repository) {
	cfg := config.Get()
	authService := auth.New(repo, cfg)
	handler := NewHandler(authService)

	authGroup := router.Group(path)
	{
		// Apply rate limiting to authentication endpoints
		authGroup.POST("/login", middleware.RateLimiter(), handler.login)
		authGroup.POST("/register", middleware.RateLimiter(), handler.register)
		authGroup.POST("/verify", handler.verifyToken)
		authGroup.POST("/logout", handler.logout)
	}
} 