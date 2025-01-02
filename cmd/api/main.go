package main

import (
	"log"

	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/database"
	"github.com/PakornBank/learn-go/internal/router"
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

	r := gin.Default()
	router.NewRouter(r, db, config).SetupRoutes()

	log.Printf("Server running on port %s\n", config.ServerPort)
	if err := r.Run(":" + config.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
