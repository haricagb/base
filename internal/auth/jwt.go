// internal/auth/jwt.go
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType distinguishes access tokens from refresh tokens.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT claims for this application.
type Claims struct {
	jwt.RegisteredClaims
	UserID       int64     `json:"user_id"`
	Role         string    `json:"role"`
	EnterpriseID int64     `json:"enterprise_id,omitempty"`
	RobotID      int64     `json:"robot_id,omitempty"`
	Type         TokenType `json:"type"`
}

// TokenInput holds the fields needed to generate a token pair.
type TokenInput struct {
	UserID       int64
	Role         string
	EnterpriseID *int64
	RobotID      *int64
}

// JWTConfig holds the settings needed by JWT operations.
type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	Issuer             string
}

// JWTManager handles token generation and validation.
type JWTManager struct {
	config JWTConfig
}

// NewJWTManager creates a JWTManager.
func NewJWTManager(config JWTConfig) *JWTManager {
	return &JWTManager{config: config}
}

// GenerateAccessToken creates a signed access token for the given user.
func (m *JWTManager) GenerateAccessToken(input TokenInput) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			Subject:   fmt.Sprintf("%d", input.UserID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.AccessTokenExpiry)),
		},
		UserID: input.UserID,
		Role:   input.Role,
		Type:   AccessToken,
	}

	if input.EnterpriseID != nil {
		claims.EnterpriseID = *input.EnterpriseID
	}
	if input.RobotID != nil {
		claims.RobotID = *input.RobotID
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.Secret))
}

// GenerateRefreshToken creates a signed refresh token for the given user.
func (m *JWTManager) GenerateRefreshToken(userID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.config.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.config.RefreshTokenExpiry)),
		},
		UserID: userID,
		Type:   RefreshToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.Secret))
}

// ValidateToken parses and validates a token string, returning the claims.
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.config.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
