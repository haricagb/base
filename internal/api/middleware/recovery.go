// internal/api/middleware/recovery.go
package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"my-application/internal/api/interceptor"
)

// Recovery returns a middleware that recovers from panics and logs the stack trace.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()
				logger.Error("panic recovered",
					slog.String("panic", fmt.Sprintf("%v", rec)),
					slog.String("stack", string(stack)),
					slog.String("method", c.Request.Method),
					slog.String("path", c.Request.URL.Path),
				)
				interceptor.Abort(c, http.StatusInternalServerError, "internal server error", nil)
			}
		}()
		c.Next()
	}
}
