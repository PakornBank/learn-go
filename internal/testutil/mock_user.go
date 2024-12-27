package testutil

import (
	"time"

	"github.com/PakornBank/learn-go/internal/model"
	"github.com/google/uuid"
)

func NewMockUser() model.User {
	return model.User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
