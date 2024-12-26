package router

import (
	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/handler"
	"github.com/PakornBank/learn-go/internal/middleware"
	"github.com/gin-gonic/gin"
)

func setupProtectedRoutes(router *gin.Engine, cfg *config.Config, authHandler *handler.AuthHandler) {
	r := router.Group("/api")
	r.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		r.GET("/profile", authHandler.GetProfile)
	}
}
