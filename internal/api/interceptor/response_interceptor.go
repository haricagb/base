// internal/api/interceptor/response_interceptor.go
package interceptor

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestIDKey is the gin context key for the request ID.
const RequestIDKey = "request_id"

// APIResponse is the standard JSON envelope for ALL API responses.
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Errors    interface{} `json:"errors"`
	RequestID string      `json:"request_id"`
	Timestamp string      `json:"timestamp"`
}

// newEnvelope constructs a standard API response envelope.
func newEnvelope(c *gin.Context, success bool, message string, data, errors interface{}) APIResponse {
	rid, _ := c.Get(RequestIDKey)
	requestID, _ := rid.(string) //nolint:errcheck // type assertion fallback to empty string is intended

	return APIResponse{
		Success:   success,
		Message:   message,
		Data:      data,
		Errors:    errors,
		RequestID: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Success sends a successful JSON response wrapped in the standard envelope.
func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, newEnvelope(c, true, http.StatusText(status), data, nil))
}

// SuccessWithMessage sends a successful JSON response with a custom message.
func SuccessWithMessage(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, newEnvelope(c, true, message, data, nil))
}

// Fail sends an error JSON response wrapped in the standard envelope.
func Fail(c *gin.Context, status int, message string, errors interface{}) {
	c.JSON(status, newEnvelope(c, false, message, nil, errors))
}

// Abort sends an error response and aborts the middleware chain.
func Abort(c *gin.Context, status int, message string, errors interface{}) {
	c.AbortWithStatusJSON(status, newEnvelope(c, false, message, nil, errors))
}

// HandleNoRoute returns a handler for unmatched routes (404).
func HandleNoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		Fail(c, http.StatusNotFound, "route not found", nil)
	}
}

// HandleNoMethod returns a handler for unsupported methods (405).
func HandleNoMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		Fail(c, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
}
