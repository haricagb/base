// internal/api/router/router.go
package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"my-application/internal/api/handler"
	"my-application/internal/api/interceptor"
	"my-application/internal/api/middleware"
	"my-application/internal/auth"
)

// Config holds middleware configuration needed by the router.
type Config struct {
	CORSConfig        middleware.CORSConfig
	RateLimitRPS      float64
	RateLimitBurst    int
	GinMode           string
	InternalAPISecret string
}

// New creates and configures the Gin engine with all middleware and routes.
func New(
	h *handler.Handler,
	authHandler *auth.Handler,
	actionsHandler *handler.ActionsHandler,
	jwtManager *auth.JWTManager,
	cfg Config,
	logger *slog.Logger,
) *gin.Engine {
	if cfg.GinMode != "" {
		gin.SetMode(cfg.GinMode)
	}

	r := gin.New()

	// Custom handlers for unmatched routes and methods.
	r.NoRoute(interceptor.HandleNoRoute())
	r.NoMethod(interceptor.HandleNoMethod())

	// Global middleware chain (outermost → innermost).
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging(logger))
	r.Use(middleware.CORS(cfg.CORSConfig))
	r.Use(middleware.RateLimit(cfg.RateLimitRPS, cfg.RateLimitBurst, logger))

	// Public routes (no auth required).
	r.GET("/health", h.Health.HealthCheck)
	r.GET("/ping", h.Health.Ping)

	// API v1 routes.
	v1 := r.Group("/api/v1")
	{
		// Public auth routes (no JWT required).
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/firebase-login", authHandler.FirebaseLogin)

			// sync-user is called by Firebase Cloud Function (server-to-server).
			syncUser := authGroup.Group("")
			syncUser.Use(middleware.InternalAuth(cfg.InternalAPISecret, logger))
			{
				syncUser.POST("/sync-user", authHandler.SyncUser)
			}
		}

		// Protected routes (JWT required).
		protected := v1.Group("")
		protected.Use(middleware.Auth(jwtManager, logger))
		{
			users := protected.Group("/users")
			{
				// mta and eta can list and create users.
				users.GET("", middleware.RequireRole("mta", "eta"), h.User.List)
				users.POST("", middleware.RequireRole("mta", "eta"), h.User.Create)

				// All authenticated roles can view a user by ID.
				users.GET("/:id", middleware.RequireRole("mta", "eta", "caregiver", "family"), h.User.GetByID)

				// mta and eta can update users.
				users.PUT("/:id", middleware.RequireRole("mta", "eta"), h.User.Update)

				// Only mta can delete users.
				users.DELETE("/:id", middleware.RequireRole("mta"), h.User.Delete)
			}
		}
	}

	// Hasura Actions webhook endpoints (called by Hasura, not clients directly).
	actions := v1.Group("/actions")
	{
		actions.POST("/login", actionsHandler.Login)
		actions.POST("/register", actionsHandler.Register)
		actions.POST("/refresh", actionsHandler.Refresh)
		actions.POST("/sync-user", actionsHandler.SyncUser)
	}

	return r
}
