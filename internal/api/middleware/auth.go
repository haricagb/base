// internal/api/middleware/auth.go
package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"my-application/internal/api/interceptor"
	"my-application/internal/auth"
)

// Context keys for authenticated user data.
const (
	ContextKeyUserID   = "auth_user_id"
	ContextKeyUserRole = "auth_user_role"
)

// Auth returns a middleware that validates JWT tokens using the auth module.
func Auth(jwtManager *auth.JWTManager, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			interceptor.Abort(c, http.StatusUnauthorized, "missing authorization header", nil)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			interceptor.Abort(c, http.StatusUnauthorized, "invalid authorization format", nil)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			interceptor.Abort(c, http.StatusUnauthorized, "empty token", nil)
			return
		}

		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			logger.Debug("token validation failed",
				slog.String("error", err.Error()),
			)
			interceptor.Abort(c, http.StatusUnauthorized, "invalid or expired token", nil)
			return
		}

		// Ensure this is an access token, not a refresh token.
		if claims.Type != auth.AccessToken {
			interceptor.Abort(c, http.StatusUnauthorized, "invalid token type", nil)
			return
		}

		// Set user info on context for downstream handlers.
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUserRole, claims.Role)

		logger.Debug("auth middleware passed",
			slog.Int64("user_id", claims.UserID),
			slog.String("role", claims.Role),
		)

		c.Next()
	}
}

// InternalAuth returns a middleware that validates server-to-server calls
// using a shared secret in the X-Internal-Secret header.
func InternalAuth(secret string, logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret == "" {
			// No secret configured — allow (dev mode).
			logger.Warn("internal auth middleware: no secret configured, allowing request")
			c.Next()
			return
		}

		provided := c.GetHeader("X-Internal-Secret")
		if provided == "" || provided != secret {
			logger.Debug("internal auth failed: invalid or missing secret")
			interceptor.Abort(c, http.StatusUnauthorized, "unauthorized: invalid internal secret", nil)
			return
		}

		c.Next()
	}
}
