package router

import (
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Register(*gin.Context)
	Login(*gin.Context)
	GetProfile(*gin.Context)
}

func setupAuthRoutes(router *gin.Engine, authHandler Handler) {
	r := router.Group("/api")
	{
		r.POST("/register", authHandler.Register)
		r.POST("/login", authHandler.Login)
	}
}
