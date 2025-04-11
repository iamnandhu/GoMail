package middleware

import (
	"GoMail/app/config"
	"GoMail/app/logic/auth"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Simple in-memory rate limiter
type rateLimiter struct {
	mu        sync.Mutex
	ipLimits  map[string][]time.Time
	maxLimit  int
	timeFrame time.Duration
}

// Create a new global rate limiter instance
var authLimiter = &rateLimiter{
	ipLimits:  make(map[string][]time.Time),
	maxLimit:  5,          // 5 attempts
	timeFrame: time.Minute, // per minute
}

// RateLimiter limits the number of authentication requests per IP address
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		authLimiter.mu.Lock()
		defer authLimiter.mu.Unlock()
		
		// Clean up old requests
		now := time.Now()
		oldestAllowed := now.Add(-authLimiter.timeFrame)
		
		// Get existing requests for this IP
		requests, exists := authLimiter.ipLimits[ip]
		if exists {
			// Filter out old requests
			var recent []time.Time
			for _, req := range requests {
				if req.After(oldestAllowed) {
					recent = append(recent, req)
				}
			}
			authLimiter.ipLimits[ip] = recent
			requests = recent
		}
		
		// Check if limit exceeded
		if len(requests) >= authLimiter.maxLimit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Try again later.",
			})
			return
		}
		
		// Add current request
		authLimiter.ipLimits[ip] = append(authLimiter.ipLimits[ip], now)
		
		c.Next()
	}
}

// VerifyAuthToken is a middleware function that verifies the auth token
func VerifyAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the auth token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token required",
			})
			return
		}
		
		// Format should be "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Expected 'Bearer <token>'",
			})
			return
		}
		
		token := parts[1]
		
		// Verify the token
		cfg := config.Get()
		authService := auth.New(nil, cfg) // We only need config for token verification
		
		claims, err := authService.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}
		
		// Store user info in context for later use
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		
		c.Next()
	}
}

// SecurityHeaders adds security-related headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Cache-Control", "no-store")
		c.Header("Pragma", "no-cache")
		
		c.Next()
	}
}

// Logger is a middleware function that logs the request
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate resolution time
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Get request details
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		// Log request details
		c.Writer.Header().Set("X-Response-Time", latency.String())
		gin.DefaultWriter.Write([]byte("[GIN] " + endTime.Format("2006/01/02 - 15:04:05") + " | " + 
			method + " | " + path + " | " + time.Since(startTime).String() + " | " + 
			c.ClientIP() + " | " + string(rune(statusCode)) + "\n"))
	}
}

// Recovery returns a middleware that recovers from any panics
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
} 