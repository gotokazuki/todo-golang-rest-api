package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggingMiddleware returns a gin middleware for request logging
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request start time
		start := time.Now()

		// Get request information
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Extract client IP
		clientIP := c.ClientIP()
		if forwardedFor := c.Request.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			// Use the first IP in the X-Forwarded-For header
			clientIP = strings.Split(forwardedFor, ",")[0]
		}

		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		// Process request
		c.Next()

		// Get response information
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log output
		logger.Info("Request completed",
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.String("latency", latency.String()),
			zap.String("client_ip", clientIP),
			zap.String("method", method),
			zap.String("user_agent", userAgent),
			zap.String("error", errorMessage),
		)
	}
}
