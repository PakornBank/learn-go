package database

import (
	"fmt"

	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDataBase(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.GetDBURL()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)

	}

	return db, nil
}