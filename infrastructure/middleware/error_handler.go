package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandlerMiddleware returns a gin middleware for error handling
func ErrorHandlerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error("Request error",
					zap.Error(err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Int("status", c.Writer.Status()),
				)
			}
		}
	}
}
