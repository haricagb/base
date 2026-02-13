// internal/domain/user.go
package domain

import "time"

// Pagination defaults.
const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// User represents the core user domain entity.
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	EnterpriseID *int64    `json:"enterprise_id,omitempty"`
	FirebaseUID  *string   `json:"firebase_uid,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserFilter holds optional query parameters for listing users.
type UserFilter struct {
	Role     string
	IsActive *bool
	Limit    int
	Offset   int
}

// Normalize applies pagination defaults and clamps values to safe bounds.
func (f *UserFilter) Normalize() {
	if f.Limit <= 0 {
		f.Limit = DefaultPageSize
	}
	if f.Limit > MaxPageSize {
		f.Limit = MaxPageSize
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}
