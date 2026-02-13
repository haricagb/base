// internal/auth/response.go
package auth

import "time"

// TokenPair holds the access and refresh tokens returned to the client.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// AuthResponse is the JSON response for login and register endpoints.
type AuthResponse struct {
	User   UserInfo  `json:"user"`
	Tokens TokenPair `json:"tokens"`
}

// UserInfo is a subset of user data safe for auth responses.
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}
