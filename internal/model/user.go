package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id" validate:"required"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-" validate:"required"`
	FullName     string    `gorm:"type:varchar(255);not null" json:"full_name" validate:"required"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
