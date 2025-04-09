package auth

import (
	"github.com/gin-gonic/gin"
)

// AddRoute adds auth routes to the given router
func AddRoute(router *gin.RouterGroup, path string) {
	handler := NewHandler()

	authGroup := router.Group(path)
	{
		authGroup.POST("/login", handler.login)
		authGroup.POST("/register", handler.register)
		authGroup.POST("/verify", handler.verifyToken)
		authGroup.POST("/logout", handler.logout)
	}
} 