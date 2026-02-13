// internal/auth/request.go
package auth

// LoginRequest is the JSON body for POST /api/v1/auth/login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterRequest is the JSON body for POST /api/v1/auth/register.
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
	FullName string `json:"full_name" binding:"required"`
}

// RefreshRequest is the JSON body for POST /api/v1/auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// SyncUserRequest is the JSON body from the Firebase Cloud Function.
type SyncUserRequest struct {
	FirebaseUID string `json:"firebase_uid" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	DisplayName string `json:"display_name"`
}

// FirebaseLoginRequest is the JSON body for POST /api/v1/auth/firebase-login.
type FirebaseLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}
