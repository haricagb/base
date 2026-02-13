// internal/api/middleware/rate_limit.go
package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"my-application/internal/api/interceptor"
)

// RateLimit returns a middleware that applies a global token-bucket rate limiter.
func RateLimit(rps float64, burst int, logger *slog.Logger) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		if !limiter.Allow() {
			logger.Warn("rate limit exceeded",
				slog.String("remote_addr", c.ClientIP()),
				slog.String("path", c.Request.URL.Path),
			)
			interceptor.Abort(c, http.StatusTooManyRequests, "rate limit exceeded", nil)
			return
		}
		c.Next()
	}
}
