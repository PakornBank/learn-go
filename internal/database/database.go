// Package database provides database connection and configuration functionality.
package database

import (
	"fmt"

	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDataBase initializes a new database connection using the provided configuration.
// It connects to a PostgreSQL database using the DBURL from the config and performs
// auto-migration for the User model.
//
// Parameters:
//   - config: A pointer to a config.Config struct containing the database configuration.
//
// Returns:
//   - *gorm.DB: A pointer to the initialized gorm.DB instance.
//   - error: An error if the connection or migration fails, otherwise nil.
func NewDataBase(config *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DBURL()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}
