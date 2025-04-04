package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

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