package auth

import (
	"GoMail/app/logic/auth"
	"GoMail/app/repository/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles authentication-related HTTP requests
type Handler struct {
	authService auth.Service
}

// NewHandler creates a new auth handler
func NewHandler(authService auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// login handles user login
func (h *Handler) login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Authenticate user and get token
	token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if err == auth.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, LoginResponse{
				Success: false,
				Error:   "Invalid email or password",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Success: false,
			Error:   "Failed to login",
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:   token,
		Success: true,
	})
}

// register handles user registration
func (h *Handler) register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create user object
	user := &models.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Register user
	err := h.authService.Register(c.Request.Context(), user, req.Password)
	if err != nil {
		if err == auth.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, RegisterResponse{
				Success: false,
				Error:   "User with this email already exists",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, RegisterResponse{
			Success: false,
			Error:   "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

// verifyToken validates a token
func (h *Handler) verifyToken(c *gin.Context) {
	var req VerifyTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify token
	_, err := h.authService.VerifyToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, VerifyTokenResponse{
			Valid: false,
			Error: "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, VerifyTokenResponse{
		Valid: true,
	})
}

// logout handles user logout
func (h *Handler) logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Revoke the token
	err := h.authService.RevokeToken(c.Request.Context(), req.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, LogoutResponse{
			Success: false,
			Error:   "Failed to logout",
		})
		return
	}
	
	c.JSON(http.StatusOK, LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	})
} 