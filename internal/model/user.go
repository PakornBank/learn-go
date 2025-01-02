// Package model contains the data structures and models used in the application.
package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
// It contains the user's unique identifier, email, password hash, full name, and timestamps for creation and updates.
//
// Fields:
//   - ID: A unique identifier for the user, generated automatically.
//   - Email: The user's email address, which must be unique and not null.
//   - PasswordHash: A hashed version of the user's password, which is required and not exposed in JSON responses.
//   - FullName: The user's full name, which is required.
//   - CreatedAt: The timestamp when the user was created, with a default value of the current timestamp.
//   - UpdatedAt: The timestamp when the user was last updated, with a default value of the current timestamp.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" validate:"required"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-" validate:"required"`
	FullName     string    `gorm:"type:varchar(255);not null" json:"full_name" validate:"required"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
