package router

import (
	"github.com/PakornBank/learn-go/internal/handler"
	"github.com/PakornBank/learn-go/internal/middleware"
	"github.com/PakornBank/learn-go/internal/repository"
	"github.com/PakornBank/learn-go/internal/service"
)

func (r *Router) setupAuthRoutes() {
	handler := handler.NewAuthHandler(service.NewAuthService(repository.NewUserRepository(r.db), r.config))

	group := r.group.Group("/auth")
	{
		group.POST("/register", handler.Register)
		group.POST("/login", handler.Login)
	}

	protected := group.Group("")
	protected.Use(middleware.AuthMiddleware(r.config.JWTSecret))
	{
		protected.GET("/profile", handler.GetProfile)
	}
}
