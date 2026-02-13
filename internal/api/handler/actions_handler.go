package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"my-application/internal/auth"
)

// ActionsHandler handles Hasura Action webhook requests.
type ActionsHandler struct {
	authService auth.Service
	logger      *slog.Logger
}

// NewActionsHandler creates an ActionsHandler.
func NewActionsHandler(authService auth.Service, logger *slog.Logger) *ActionsHandler {
	return &ActionsHandler{authService: authService, logger: logger}
}

// hasuraActionPayload is the standard envelope Hasura sends for synchronous Actions.
type hasuraActionPayload struct {
	Action struct {
		Name string `json:"name"`
	} `json:"action"`
	Input       map[string]interface{} `json:"input"`
	SessionVars map[string]string      `json:"session_variables"`
}

// Login handles the Hasura "login" action.
func (h *ActionsHandler) Login(c *gin.Context) {
	var payload hasuraActionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	req := auth.LoginRequest{
		Email:    payload.Input["email"].(string),
		Password: payload.Input["password"].(string),
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("action login failed", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Register handles the Hasura "register" action.
func (h *ActionsHandler) Register(c *gin.Context) {
	var payload hasuraActionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	req := auth.RegisterRequest{
		Username: payload.Input["username"].(string),
		Email:    payload.Input["email"].(string),
		Password: payload.Input["password"].(string),
		FullName: payload.Input["full_name"].(string),
	}

	resp, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("action register failed", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Refresh handles the Hasura "refreshToken" action.
func (h *ActionsHandler) Refresh(c *gin.Context) {
	var payload hasuraActionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	req := auth.RefreshRequest{
		RefreshToken: payload.Input["refresh_token"].(string),
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), req)
	if err != nil {
		h.logger.Warn("action refresh failed", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// SyncUser handles the Hasura "syncUser" action.
func (h *ActionsHandler) SyncUser(c *gin.Context) {
	var payload hasuraActionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	req := auth.SyncUserRequest{
		FirebaseUID: payload.Input["firebase_uid"].(string),
		Email:       payload.Input["email"].(string),
	}
	if dn, ok := payload.Input["display_name"].(string); ok {
		req.DisplayName = dn
	}

	resp, err := h.authService.SyncUser(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("action sync-user failed", slog.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
