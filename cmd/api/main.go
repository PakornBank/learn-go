package main

import (
	"log"

	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/database"
	"github.com/PakornBank/learn-go/internal/handler"
	"github.com/PakornBank/learn-go/internal/middleware"
	"github.com/PakornBank/learn-go/internal/repository"
	"github.com/PakornBank/learn-go/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewDataBase(config)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, config)
	authHandler := handler.NewAuthHandler(authService)

	r := gin.Default()

	r.POST("/api/register", authHandler.Register)
	r.POST("/api/login", authHandler.Login)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(config.JWTSecret))
	{
		protected.GET("/profile", authHandler.GetProfile)
	}

	log.Printf("Server running on port %s\n", config.ServerPort)
	if err := r.Run(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
