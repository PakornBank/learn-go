package router

import (
	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.Default()

	setupAuthRoutes(r, authHandler)
	setupProtectedRoutes(r, cfg, authHandler)

	return r
}
