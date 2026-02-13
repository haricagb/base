// internal/api/middleware/request_id.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"my-application/internal/api/interceptor"
)

// RequestID generates or extracts a request ID for tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.New().String()
		}

		c.Set(interceptor.RequestIDKey, id)
		c.Header("X-Request-ID", id)

		c.Next()
	}
}
