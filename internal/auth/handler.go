// internal/auth/handler.go
package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"my-application/internal/api/interceptor"
	"my-application/internal/domain"
	"my-application/pkg/logger"
)

// Handler handles authentication HTTP requests.
type Handler struct {
	authService Service
	logger      *slog.Logger
}

// NewHandler creates an auth Handler.
func NewHandler(authService Service, logger *slog.Logger) *Handler {
	return &Handler{authService: authService, logger: logger}
}

// Login handles POST /api/v1/auth/login
func (h *Handler) Login(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		interceptor.Fail(c, http.StatusBadRequest, "invalid request: "+err.Error(), nil)
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		log.Warn("login failed", slog.String("email", req.Email), slog.String("error", err.Error()))
		respondAuthError(c, err)
		return
	}

	interceptor.Success(c, http.StatusOK, resp)
}

// Register handles POST /api/v1/auth/register
func (h *Handler) Register(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		interceptor.Fail(c, http.StatusBadRequest, "invalid request: "+err.Error(), nil)
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		log.Error("registration failed", slog.String("error", err.Error()))
		respondAuthError(c, err)
		return
	}

	interceptor.Success(c, http.StatusCreated, resp)
}

// Refresh handles POST /api/v1/auth/refresh
func (h *Handler) Refresh(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		interceptor.Fail(c, http.StatusBadRequest, "invalid request: "+err.Error(), nil)
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		log.Warn("token refresh failed", slog.String("error", err.Error()))
		respondAuthError(c, err)
		return
	}

	interceptor.Success(c, http.StatusOK, tokens)
}

// SyncUser handles POST /api/auth/sync-user (called by Firebase Cloud Function).
func (h *Handler) SyncUser(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var req SyncUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		interceptor.Fail(c, http.StatusBadRequest, "invalid request: "+err.Error(), nil)
		return
	}

	resp, err := h.authService.SyncUser(c.Request.Context(), req)
	if err != nil {
		log.Error("user sync failed",
			slog.String("firebase_uid", req.FirebaseUID),
			slog.String("error", err.Error()),
		)
		respondAuthError(c, err)
		return
	}

	interceptor.Success(c, http.StatusOK, resp)
}

// FirebaseLogin handles POST /api/v1/auth/firebase-login
func (h *Handler) FirebaseLogin(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var req FirebaseLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		interceptor.Fail(c, http.StatusBadRequest, "invalid request: "+err.Error(), nil)
		return
	}

	resp, err := h.authService.FirebaseLogin(c.Request.Context(), req)
	if err != nil {
		log.Warn("firebase login failed", slog.String("error", err.Error()))
		respondAuthError(c, err)
		return
	}

	interceptor.Success(c, http.StatusOK, resp)
}

// respondAuthError maps domain errors to HTTP responses.
func respondAuthError(c *gin.Context, err error) {
	appErr, ok := err.(*domain.AppError)
	if ok {
		status := domain.HTTPStatusFromError(appErr.Err)
		interceptor.Fail(c, status, appErr.Message, appErr.Details)
		return
	}
	interceptor.Fail(c, http.StatusInternalServerError, "internal server error", nil)
}
