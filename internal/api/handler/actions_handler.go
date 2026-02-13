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

// inputString safely extracts a string value from the Hasura action input map.
func (p *hasuraActionPayload) inputString(key string) (string, bool) {
	v, ok := p.Input[key].(string)
	return v, ok
}

// Login handles the Hasura "login" action.
func (h *ActionsHandler) Login(c *gin.Context) {
	var payload hasuraActionPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid payload"})
		return
	}

	email, ok := payload.inputString("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required and must be a string"})
		return
	}
	password, ok := payload.inputString("password")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is required and must be a string"})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), auth.LoginRequest{
		Email:    email,
		Password: password,
	})
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

	username, ok := payload.inputString("username")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username is required and must be a string"})
		return
	}
	email, ok := payload.inputString("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required and must be a string"})
		return
	}
	password, ok := payload.inputString("password")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "password is required and must be a string"})
		return
	}
	fullName, ok := payload.inputString("full_name")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "full_name is required and must be a string"})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), auth.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
		FullName: fullName,
	})
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

	refreshToken, ok := payload.inputString("refresh_token")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "refresh_token is required and must be a string"})
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), auth.RefreshRequest{
		RefreshToken: refreshToken,
	})
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

	firebaseUID, ok := payload.inputString("firebase_uid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "firebase_uid is required and must be a string"})
		return
	}
	email, ok := payload.inputString("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "email is required and must be a string"})
		return
	}

	req := auth.SyncUserRequest{
		FirebaseUID: firebaseUID,
		Email:       email,
	}
	if dn, ok := payload.inputString("display_name"); ok {
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
