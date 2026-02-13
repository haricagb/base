// internal/api/middleware/rbac.go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"my-application/internal/api/interceptor"
)

// RequireRole returns a middleware that checks if the authenticated user
// has one of the allowed roles. Must be used after the Auth middleware.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyUserRole)
		if !exists {
			interceptor.Abort(c, http.StatusForbidden, "access denied: no role found", nil)
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			interceptor.Abort(c, http.StatusForbidden, "access denied: invalid role", nil)
			return
		}

		if _, allowed := roleSet[roleStr]; !allowed {
			interceptor.Abort(c, http.StatusForbidden, "access denied: insufficient permissions", nil)
			return
		}

		c.Next()
	}
}
