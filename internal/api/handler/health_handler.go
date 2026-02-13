// internal/api/handler/health_handler.go
package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"my-application/internal/api/interceptor"
	"my-application/pkg/database"
)

// HealthHandler handles health check and test endpoints.
type HealthHandler struct {
	dbPool *pgxpool.Pool
	logger *slog.Logger
}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler(dbPool *pgxpool.Pool, logger *slog.Logger) *HealthHandler {
	return &HealthHandler{dbPool: dbPool, logger: logger}
}

// HealthCheck returns the overall health of the service including database connectivity.
// GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	dbStatus := "up"
	if err := database.HealthCheck(c.Request.Context(), h.dbPool); err != nil {
		dbStatus = "down"
		h.logger.Error("database health check failed", slog.String("error", err.Error()))
	}

	status := http.StatusOK
	serviceStatus := "healthy"
	if dbStatus != "up" {
		status = http.StatusServiceUnavailable
		serviceStatus = "degraded"
	}

	interceptor.Success(c, status, gin.H{
		"status": serviceStatus,
		"checks": gin.H{
			"database": dbStatus,
		},
	})
}

// Ping is a lightweight liveness probe.
// GET /ping
func (h *HealthHandler) Ping(c *gin.Context) {
	interceptor.SuccessWithMessage(c, http.StatusOK, "pong", nil)
}
