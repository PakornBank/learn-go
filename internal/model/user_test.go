package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validation(t *testing.T) {
	validate := validator.New()
	testUUID := uuid.New()
	now := time.Now()
	const (
		testEmail    = "test@example.com"
		testPassword = "hashedpassword"
		testFullName = "Test User"
	)

	tests := []struct {
		name        string
		user        User
		wantErr     bool
		errContains string
	}{
		{
			name: "valid user",
			user: User{
				ID:           testUUID,
				Email:        testEmail,
				FullName:     testFullName,
				PasswordHash: testPassword,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr: false,
		},
		{
			name: "missing email",
			user: User{
				ID:           testUUID,
				Email:        "",
				FullName:     testFullName,
				PasswordHash: testPassword,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr:     true,
			errContains: "Email",
		},
		{
			name: "invalid email",
			user: User{
				ID:           testUUID,
				Email:        "invalid-email",
				FullName:     testFullName,
				PasswordHash: testPassword,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr:     true,
			errContains: "Email",
		},
		{
			name: "missing password hash",
			user: User{
				ID:           testUUID,
				Email:        testEmail,
				FullName:     testFullName,
				PasswordHash: "",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr:     true,
			errContains: "PasswordHash",
		},
		{
			name: "missing full name",
			user: User{
				ID:           testUUID,
				Email:        testEmail,
				FullName:     "",
				PasswordHash: testPassword,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr:     true,
			errContains: "FullName",
		},
		{
			name: "zero value UUID",
			user: User{
				ID:           uuid.UUID{},
				Email:        testEmail,
				FullName:     testFullName,
				PasswordHash: testPassword,
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr:     true,
			errContains: "ID",
		},
		{
			name: "zero value created time",
			user: User{
				ID:           testUUID,
				Email:        testEmail,
				FullName:     testFullName,
				PasswordHash: testPassword,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
			wantErr: false,
		},
		{
			name: "missing optional fields",
			user: User{
				ID:           testUUID,
				Email:        testEmail,
				FullName:     testFullName,
				PasswordHash: testPassword,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_JSONSerialization(t *testing.T) {
	user := User{
		ID:           uuid.New(),
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		FullName:     "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("password hash should not be serialized", func(t *testing.T) {
		jsonData, err := json.Marshal(user)
		assert.NoError(t, err)

		var unmarshaled User
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)

		assert.Empty(t, unmarshaled.PasswordHash)
		assert.Equal(t, user.ID, unmarshaled.ID)
		assert.Equal(t, user.Email, unmarshaled.Email)
		assert.Equal(t, user.FullName, unmarshaled.FullName)
		assert.Equal(t, user.CreatedAt.Truncate(time.Microsecond), unmarshaled.CreatedAt)
		assert.Equal(t, user.UpdatedAt.Truncate(time.Microsecond), unmarshaled.UpdatedAt)
	})
}
