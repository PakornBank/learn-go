package router

import (
	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/handler"
	"github.com/PakornBank/learn-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

// setupProtectedRoutes adds authenticated API routes to the router.
func setupProtectedRoutes(r *gin.Engine, cfg *config.Config, auth *handler.AuthHandler) {
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/profile", auth.GetProfile)
	}
}
