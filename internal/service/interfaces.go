// internal/service/interfaces.go
package service

import (
	"context"

	"my-application/internal/domain"
)

// UserService defines business operations for Users.
type UserService interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id int64) error
}
