// cmd/api/main.go
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"my-application/config"
	"my-application/internal/api/handler"
	"my-application/internal/api/middleware"
	"my-application/internal/api/router"
	"my-application/internal/auth"
	"my-application/internal/repository/postgres"
	"my-application/internal/service"
	"my-application/pkg/database"
	fbclient "my-application/pkg/firebase"
	"my-application/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Configuration.
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	cfg, err := config.Load("config", env)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// 2. Logger.
	log := logger.Setup(cfg.Log.Level, cfg.Log.Format, os.Stdout)
	log.Info("starting application",
		slog.String("env", env),
		slog.Int("port", cfg.Server.Port),
	)

	// 3. Gin mode.
	ginMode := gin.ReleaseMode
	switch env {
	case "dev", "development":
		ginMode = gin.DebugMode
	case "test":
		ginMode = gin.TestMode
	}

	// 4. Firebase (optional — skip if not configured).
	ctx := context.Background()
	var firebaseVerifier auth.FirebaseVerifier
	if cfg.Firebase.ProjectID != "" && cfg.Firebase.CredentialsFile != "" {
		fbClient, fbErr := fbclient.NewClient(ctx, fbclient.Config{
			ProjectID:       cfg.Firebase.ProjectID,
			CredentialsFile: cfg.Firebase.CredentialsFile,
		}, log)
		if fbErr != nil {
			return fmt.Errorf("initializing firebase: %w", fbErr)
		}
		firebaseVerifier = &firebaseAdapter{client: fbClient}
		log.Info("Firebase authentication enabled")
	} else {
		log.Warn("Firebase not configured — firebase-login endpoint will return an error")
	}

	// 5. Database.
	dbPool, err := database.NewPostgresPool(ctx, database.PostgresConfig{
		DSN:             cfg.Database.DSN(),
		MaxConns:        cfg.Database.MaxConns,
		MinConns:        cfg.Database.MinConns,
		MaxConnLifetime: cfg.Database.MaxConnLifetime,
		MaxConnIdleTime: cfg.Database.MaxConnIdleTime,
	}, log)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer dbPool.Close()

	// 6. Repository layer.
	userRepo := postgres.NewUserPostgres(dbPool, log)

	// 7. Service layer.
	userSvc := service.NewUserService(userRepo, log)

	// 8. Handler layer.
	h := handler.NewHandler(userSvc, dbPool, log)

	// 9. Auth module.
	jwtManager := auth.NewJWTManager(auth.JWTConfig{
		Secret:             cfg.JWT.Secret,
		AccessTokenExpiry:  cfg.JWT.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.JWT.RefreshTokenExpiry,
		Issuer:             cfg.JWT.Issuer,
	})
	authSvc := auth.NewService(userRepo, jwtManager, firebaseVerifier, log)
	authHandler := auth.NewHandler(authSvc, log)
	actionsHandler := handler.NewActionsHandler(authSvc, log)

	// 10. Router.
	r := router.New(h, authHandler, actionsHandler, jwtManager, router.Config{
		CORSConfig: middleware.CORSConfig{
			AllowedOrigins:   cfg.CORS.AllowedOrigins,
			AllowedMethods:   cfg.CORS.AllowedMethods,
			AllowedHeaders:   cfg.CORS.AllowedHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           cfg.CORS.MaxAge,
		},
		RateLimitRPS:      cfg.RateLimit.RequestsPerSecond,
		RateLimitBurst:    cfg.RateLimit.Burst,
		GinMode:           ginMode,
		InternalAPISecret: os.Getenv("INTERNAL_API_SECRET"),
	}, log)

	// 11. HTTP Server.
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 12. Graceful Shutdown.
	go func() {
		log.Info("HTTP server listening", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("shutdown signal received", slog.String("signal", sig.String()))

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	log.Info("server shutdown completed gracefully")
	return nil
}

// firebaseAdapter adapts pkg/firebase.Client to the auth.FirebaseVerifier interface.
type firebaseAdapter struct {
	client *fbclient.Client
}

func (a *firebaseAdapter) VerifyIDToken(ctx context.Context, idToken string) (*auth.VerifiedFirebaseUser, error) {
	user, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return &auth.VerifiedFirebaseUser{
		UID:         user.UID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}, nil
}
