// internal/auth/service.go
package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"my-application/internal/domain"
	"my-application/internal/repository"
)

// FirebaseVerifier abstracts Firebase ID token verification.
type FirebaseVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*VerifiedFirebaseUser, error)
}

// VerifiedFirebaseUser holds user info extracted from a verified Firebase token.
type VerifiedFirebaseUser struct {
	UID         string
	Email       string
	DisplayName string
}

// Service defines the authentication business operations.
type Service interface {
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	RefreshToken(ctx context.Context, req RefreshRequest) (*TokenPair, error)
	SyncUser(ctx context.Context, req SyncUserRequest) (*AuthResponse, error)
	FirebaseLogin(ctx context.Context, req FirebaseLoginRequest) (*AuthResponse, error)
}

// Compile-time interface check.
var _ Service = (*authService)(nil)

type authService struct {
	userRepo         repository.UserRepository
	jwtManager       *JWTManager
	firebaseVerifier FirebaseVerifier
	logger           *slog.Logger
}

// NewService creates a new auth Service.
// firebaseVerifier can be nil if Firebase is not configured.
func NewService(
	userRepo repository.UserRepository,
	jwtManager *JWTManager,
	firebaseVerifier FirebaseVerifier,
	logger *slog.Logger,
) Service {
	return &authService{
		userRepo:         userRepo,
		jwtManager:       jwtManager,
		firebaseVerifier: firebaseVerifier,
		logger:           logger,
	}
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// 1. Find user by email.
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if err != nil {
		// Map "not found" to "unauthorized" so we don't leak whether emails exist.
		var appErr *domain.AppError
		if errors.As(err, &appErr) && errors.Is(appErr.Err, domain.ErrNotFound) {
			return nil, domain.NewAppError(domain.ErrUnauthorized, "invalid email or password")
		}
		return nil, err
	}

	// 2. Check if user is active.
	if !user.IsActive {
		return nil, domain.NewAppError(domain.ErrUnauthorized, "account is deactivated")
	}

	// 3. Verify password.
	if err := CheckPassword(req.Password, user.PasswordHash); err != nil {
		return nil, domain.NewAppError(domain.ErrUnauthorized, "invalid email or password")
	}

	// 4. Generate tokens.
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		s.logger.Error("failed to generate tokens", slog.String("error", err.Error()))
		return nil, domain.NewAppError(domain.ErrInternal, "failed to generate tokens")
	}

	return &AuthResponse{
		User:   toUserInfo(user),
		Tokens: *tokens,
	}, nil
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// 1. Hash password.
	hash, err := HashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", slog.String("error", err.Error()))
		return nil, domain.NewAppError(domain.ErrInternal, "failed to process registration")
	}

	// 2. Create domain user.
	user := &domain.User{
		Username:     strings.TrimSpace(req.Username),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: hash,
		FullName:     strings.TrimSpace(req.FullName),
		Role:         "caregiver", // Default role for new SONA registrations.
		IsActive:     true,
	}

	// 3. Persist (the repository handles unique constraint violations).
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 4. Generate tokens.
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		s.logger.Error("failed to generate tokens", slog.String("error", err.Error()))
		return nil, domain.NewAppError(domain.ErrInternal, "failed to generate tokens")
	}

	return &AuthResponse{
		User:   toUserInfo(user),
		Tokens: *tokens,
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req RefreshRequest) (*TokenPair, error) {
	// 1. Validate the refresh token.
	claims, err := s.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, domain.NewAppError(domain.ErrUnauthorized, "invalid or expired refresh token")
	}

	// 2. Ensure it is actually a refresh token.
	if claims.Type != RefreshToken {
		return nil, domain.NewAppError(domain.ErrUnauthorized, "token is not a refresh token")
	}

	// 3. Verify the user still exists and is active.
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, domain.NewAppError(domain.ErrUnauthorized, "user not found")
	}
	if !user.IsActive {
		return nil, domain.NewAppError(domain.ErrUnauthorized, "account is deactivated")
	}

	// 4. Generate a new token pair.
	return s.generateTokenPair(user)
}

// SyncUser handles user creation/lookup from Firebase Cloud Function.
// If a user with the given firebase_uid already exists, return tokens for them.
// Otherwise create a new user with default role.
func (s *authService) SyncUser(ctx context.Context, req SyncUserRequest) (*AuthResponse, error) {
	return s.findOrCreateFirebaseUser(ctx, req.FirebaseUID, req.Email, req.DisplayName)
}

// FirebaseLogin verifies a Firebase ID token and returns custom JWT tokens.
// If the Firebase user doesn't exist in the local DB, they are created.
func (s *authService) FirebaseLogin(ctx context.Context, req FirebaseLoginRequest) (*AuthResponse, error) {
	if s.firebaseVerifier == nil {
		return nil, domain.NewAppError(domain.ErrInternal, "firebase authentication is not configured")
	}

	// 1. Verify the Firebase ID token.
	fbUser, err := s.firebaseVerifier.VerifyIDToken(ctx, req.IDToken)
	if err != nil {
		s.logger.Warn("firebase token verification failed", slog.String("error", err.Error()))
		return nil, domain.NewAppError(domain.ErrUnauthorized, "invalid or expired firebase token")
	}

	// 2. Find or create the local user.
	return s.findOrCreateFirebaseUser(ctx, fbUser.UID, fbUser.Email, fbUser.DisplayName)
}

// findOrCreateFirebaseUser looks up a user by firebase_uid, or creates one if not found.
func (s *authService) findOrCreateFirebaseUser(ctx context.Context, firebaseUID, email, displayName string) (*AuthResponse, error) {
	// 1. Check if user already exists by firebase_uid.
	user, err := s.userRepo.GetByFirebaseUID(ctx, firebaseUID)
	if err == nil {
		// User exists — check active, generate tokens and return.
		if !user.IsActive {
			return nil, domain.NewAppError(domain.ErrUnauthorized, "account is deactivated")
		}
		tokens, tokenErr := s.generateTokenPair(user)
		if tokenErr != nil {
			return nil, domain.NewAppError(domain.ErrInternal, "failed to generate tokens")
		}
		return &AuthResponse{
			User:   toUserInfo(user),
			Tokens: *tokens,
		}, nil
	}

	// 2. User doesn't exist — create a new one.
	email = strings.ToLower(strings.TrimSpace(email))
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		displayName = email
	}

	// Generate a username from the email prefix.
	username := strings.Split(email, "@")[0]

	user = &domain.User{
		Username:    username,
		Email:       email,
		FullName:    displayName,
		FirebaseUID: &firebaseUID,
		Role:        "caregiver", // Default SONA role for Firebase signups.
		IsActive:    true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		// If username conflict, append part of firebase_uid to make it unique.
		var appErr *domain.AppError
		if errors.As(err, &appErr) && errors.Is(appErr.Err, domain.ErrAlreadyExists) {
			user.Username = username + "_" + firebaseUID[:8]
			if createErr := s.userRepo.Create(ctx, user); createErr != nil {
				return nil, createErr
			}
		} else {
			return nil, err
		}
	}

	// 3. Generate tokens.
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		s.logger.Error("failed to generate tokens for firebase user", slog.String("error", err.Error()))
		return nil, domain.NewAppError(domain.ErrInternal, "failed to generate tokens")
	}

	return &AuthResponse{
		User:   toUserInfo(user),
		Tokens: *tokens,
	}, nil
}

func (s *authService) generateTokenPair(user *domain.User) (*TokenPair, error) {
	input := TokenInput{
		UserID:       user.ID,
		Role:         user.Role,
		EnterpriseID: user.EnterpriseID,
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(input)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.jwtManager.config.AccessTokenExpiry),
	}, nil
}

func toUserInfo(u *domain.User) UserInfo {
	return UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		FullName: u.FullName,
		Role:     u.Role,
	}
}
