// internal/api/handler/user_handler.go
package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"my-application/internal/api/interceptor"
	"my-application/internal/api/request"
	"my-application/internal/api/response"
	"my-application/internal/domain"
	"my-application/internal/service"
	"my-application/pkg/logger"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	userService service.UserService
	logger      *slog.Logger
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(userService service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

// List handles GET /api/v1/users
func (h *UserHandler) List(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	filter := domain.UserFilter{
		Role: c.Query("role"),
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = v
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if v, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = v
		}
	}
	if activeStr := c.Query("is_active"); activeStr != "" {
		v := activeStr == "true"
		filter.IsActive = &v
	}

	users, total, err := h.userService.ListUsers(c.Request.Context(), filter)
	if err != nil {
		log.Error("failed to list users", slog.String("error", err.Error()))
		respondError(c, err)
		return
	}

	userResponses := make([]response.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = toUserResponse(u)
	}

	interceptor.Success(c, http.StatusOK, response.UserListResponse{
		Users:  userResponses,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	})
}

// GetByID handles GET /api/v1/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, domain.NewAppError(domain.ErrInvalidInput, "invalid user ID"))
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		respondError(c, err)
		return
	}

	interceptor.Success(c, http.StatusOK, toUserResponse(*user))
}

// Create handles POST /api/v1/users
func (h *UserHandler) Create(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, domain.NewAppError(domain.ErrInvalidInput, "invalid JSON: "+err.Error()))
		return
	}

	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     req.Role,
	}

	if err := h.userService.CreateUser(c.Request.Context(), user); err != nil {
		log.Error("failed to create user", slog.String("error", err.Error()))
		respondError(c, err)
		return
	}

	interceptor.Success(c, http.StatusCreated, toUserResponse(*user))
}

// Update handles PUT /api/v1/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, domain.NewAppError(domain.ErrInvalidInput, "invalid user ID"))
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, domain.NewAppError(domain.ErrInvalidInput, "invalid JSON: "+err.Error()))
		return
	}

	user := &domain.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Role:     req.Role,
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := h.userService.UpdateUser(c.Request.Context(), user); err != nil {
		log.Error("failed to update user", slog.String("error", err.Error()))
		respondError(c, err)
		return
	}

	interceptor.Success(c, http.StatusOK, toUserResponse(*user))
}

// Delete handles DELETE /api/v1/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	log := logger.FromContext(c.Request.Context())

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, domain.NewAppError(domain.ErrInvalidInput, "invalid user ID"))
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		log.Error("failed to delete user", slog.String("error", err.Error()))
		respondError(c, err)
		return
	}

	interceptor.SuccessWithMessage(c, http.StatusOK, "user deleted successfully", nil)
}

func toUserResponse(u domain.User) response.UserResponse {
	return response.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
