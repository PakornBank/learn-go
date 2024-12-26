package router

import (
	"github.com/PakornBank/learn-go/internal/handler"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(router *gin.Engine, authHandler *handler.AuthHandler) {
	r := router.Group("/api")
	{
		r.POST("/register", authHandler.Register)
		r.POST("/login", authHandler.Login)
	}
}
