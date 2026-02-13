// pkg/firebase/firebase.go
package firebase

import (
	"context"
	"fmt"
	"log/slog"

	firebase "firebase.google.com/go/v4"
	fbauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// Client wraps the Firebase Auth client for token verification.
type Client struct {
	auth   *fbauth.Client
	logger *slog.Logger
}

// Config holds Firebase initialization settings.
type Config struct {
	ProjectID       string
	CredentialsFile string
}

// NewClient initializes a Firebase App and returns a Client for auth operations.
func NewClient(ctx context.Context, cfg Config, logger *slog.Logger) (*Client, error) {
	var opts []option.ClientOption
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}

	app, err := firebase.NewApp(ctx, &firebase.Config{
		ProjectID: cfg.ProjectID,
	}, opts...)
	if err != nil {
		return nil, fmt.Errorf("initializing firebase app: %w", err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("initializing firebase auth client: %w", err)
	}

	logger.Info("Firebase Auth client initialized", slog.String("project_id", cfg.ProjectID))

	return &Client{auth: authClient, logger: logger}, nil
}

// VerifiedUser holds the user info extracted from a verified Firebase ID token.
type VerifiedUser struct {
	UID         string
	Email       string
	DisplayName string
}

// VerifyIDToken validates a Firebase ID token and returns the user info.
func (c *Client) VerifyIDToken(ctx context.Context, idToken string) (*VerifiedUser, error) {
	token, err := c.auth.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("verifying firebase id token: %w", err)
	}

	email, _ := token.Claims["email"].(string) //nolint:errcheck // claims may be absent; empty string is acceptable
	name, _ := token.Claims["name"].(string)   //nolint:errcheck // claims may be absent; empty string is acceptable

	return &VerifiedUser{
		UID:         token.UID,
		Email:       email,
		DisplayName: name,
	}, nil
}
