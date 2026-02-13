// internal/api/request/user_request.go
package request

// CreateUserRequest is the JSON body for creating a user.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}

// UpdateUserRequest is the JSON body for updating a user.
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}
