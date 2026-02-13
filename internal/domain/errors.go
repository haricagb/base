// internal/domain/errors.go
package domain

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors for the application. Every layer returns these;
// the handler/interceptor maps them to HTTP status codes.
var (
	ErrNotFound          = errors.New("resource not found")
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternal          = errors.New("internal server error")
	ErrDatabaseOperation = errors.New("database operation failed")
)

// AppError wraps a sentinel error with a contextual message and optional field-level details.
type AppError struct {
	Err     error             // The sentinel error (e.g., ErrNotFound)
	Message string            // Human-readable message
	Details map[string]string // Optional field-level validation errors
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Err.Error(), e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError creates a new AppError wrapping a sentinel.
func NewAppError(err error, message string) *AppError {
	return &AppError{Err: err, Message: message}
}

// NewValidationError creates an AppError with field-level details.
func NewValidationError(message string, details map[string]string) *AppError {
	return &AppError{Err: ErrInvalidInput, Message: message, Details: details}
}

// HTTPStatusFromError maps domain errors to HTTP status codes.
func HTTPStatusFromError(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrAlreadyExists):
		return http.StatusConflict
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
