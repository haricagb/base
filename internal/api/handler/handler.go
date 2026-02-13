// internal/api/handler/handler.go
package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"my-application/internal/api/interceptor"
	"my-application/internal/domain"
	"my-application/internal/service"
)

// Handler aggregates all route handlers and shared dependencies.
type Handler struct {
	Health *HealthHandler
	User   *UserHandler
	logger *slog.Logger
}

// NewHandler creates a Handler with all sub-handlers wired up.
func NewHandler(
	userService service.UserService,
	dbPool *pgxpool.Pool,
	logger *slog.Logger,
) *Handler {
	return &Handler{
		Health: NewHealthHandler(dbPool, logger),
		User:   NewUserHandler(userService, logger),
		logger: logger,
	}
}

// respondError writes an error response mapped from a domain error.
func respondError(c *gin.Context, err error) {
	appErr, ok := err.(*domain.AppError)
	if ok {
		status := domain.HTTPStatusFromError(appErr.Err)
		interceptor.Fail(c, status, appErr.Message, appErr.Details)
		return
	}
	interceptor.Fail(c, http.StatusInternalServerError, "internal server error", nil)
}
