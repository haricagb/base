// internal/repository/interfaces.go
package repository

import (
	"context"

	"my-application/internal/domain"
)

// UserRepository defines the data access contract for User entities.
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)
	List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int64, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
}
