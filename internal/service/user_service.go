// internal/service/user_service.go
package service

import (
	"context"
	"log/slog"
	"strings"

	"my-application/internal/domain"
	"my-application/internal/repository"
)

// Compile-time interface check.
var _ UserService = (*userService)(nil)

type userService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.UserRepository, logger *slog.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *userService) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.NewAppError(domain.ErrInvalidInput, "user ID must be positive")
	}
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return s.userRepo.List(ctx, filter)
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) error {
	if err := s.validateUser(user); err != nil {
		return err
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Username = strings.TrimSpace(user.Username)
	user.IsActive = true

	if user.Role == "" {
		user.Role = "caregiver"
	}

	return s.userRepo.Create(ctx, user)
}

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.ID <= 0 {
		return domain.NewAppError(domain.ErrInvalidInput, "user ID must be positive")
	}
	if err := s.validateUser(user); err != nil {
		return err
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Username = strings.TrimSpace(user.Username)

	return s.userRepo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.NewAppError(domain.ErrInvalidInput, "user ID must be positive")
	}
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) validateUser(user *domain.User) error {
	details := make(map[string]string)

	if strings.TrimSpace(user.Username) == "" {
		details["username"] = "username is required"
	} else if len(user.Username) < 3 || len(user.Username) > 50 {
		details["username"] = "username must be between 3 and 50 characters"
	}

	if strings.TrimSpace(user.Email) == "" {
		details["email"] = "email is required"
	} else if !strings.Contains(user.Email, "@") {
		details["email"] = "email must be a valid email address"
	}

	if strings.TrimSpace(user.FullName) == "" {
		details["full_name"] = "full name is required"
	}

	validRoles := map[string]bool{"mta": true, "eta": true, "caregiver": true, "family": true, "robot": true}
	if user.Role != "" && !validRoles[user.Role] {
		details["role"] = "role must be one of: mta, eta, caregiver, family, robot"
	}

	if len(details) > 0 {
		return domain.NewValidationError("validation failed", details)
	}
	return nil
}
