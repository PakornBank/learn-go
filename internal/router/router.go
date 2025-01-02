package router

import (
	"github.com/PakornBank/learn-go/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	group  *gin.RouterGroup
	db     *gorm.DB
	config *config.Config
}

func NewRouter(r *gin.Engine, db *gorm.DB, config *config.Config) *Router {
	return &Router{
		group:  r.Group("/api"),
		db:     db,
		config: config,
	}
}

func (r *Router) SetupRoutes() {
	r.setupAuthRoutes()
}
