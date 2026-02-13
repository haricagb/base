// internal/api/middleware/logging.go
package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"my-application/internal/api/interceptor"
	"my-application/pkg/logger"
)

// Logging returns a middleware that logs each HTTP request with duration and status.
func Logging(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		rid, _ := c.Get(interceptor.RequestIDKey)
		requestID, _ := rid.(string) //nolint:errcheck // type assertion fallback to empty string is intended

		reqLogger := log.With(
			slog.String("request_id", requestID),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.ClientIP()),
		)

		ctx := logger.WithContext(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		reqLogger.Info("request completed",
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
			slog.String("user_agent", c.Request.UserAgent()),
		)
	}
}
