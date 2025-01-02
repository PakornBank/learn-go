// Package config provides configuration management for the application.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the configuration values for the application.
// It includes database connection details, server port, JWT secret, and token expiry duration.
type Config struct {
	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBPort         string
	ServerPort     string
	JWTSecret      string
	TokenExpiryDur time.Duration
}

// LoadConfig loads the configuration from environment variables and returns a Config struct.
// It first attempts to load environment variables from a .env file using godotenv.
// If the .env file does not exist, it continues without error.
// If the .env file exists but cannot be loaded, it returns an error.
//
// The following environment variables are used to populate the Config struct:
// - DB_HOST: Database host (default: "localhost")
// - DB_USER: Database user (default: "postgres")
// - DB_PASSWORD: Database password (default: "")
// - DB_NAME: Database name (default: "go_auth_db")
// - DB_PORT: Database port (default: "5432")
// - SERVER_PORT: Server port (default: "8080")
// - JWT_SECRET: JWT secret key (default: "your-secret-key")
//
// If the JWT_SECRET environment variable is not set (i.e., it is "your-secret-key"),
// the function returns an error indicating that the JWT secret must be set.
//
// Returns a pointer to a Config struct and an error, if any.
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "go_auth_db"),
		DBPort:         getEnv("DB_PORT", "5432"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiryDur: 24 * time.Hour,
	}

	if config.JWTSecret == "your-secret-key" {
		return nil, errors.New("jwt secret must be set in environment")
	}

	return config, nil
}

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is present in the environment, the function returns its value.
// Otherwise, it returns the specified defaultValue.
//
// Parameters:
//   - key: The name of the environment variable to look up.
//   - defaultValue: The value to return if the environment variable is not set.
//
// Returns:
//
//	The value of the environment variable if it exists, otherwise defaultValue.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// DBURL constructs and returns the database connection URL string
// based on the configuration fields of the Config struct.
// The returned URL includes the host, user, password, database name,
// port, and disables SSL mode.
func (c *Config) DBURL() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort,
	)
}
