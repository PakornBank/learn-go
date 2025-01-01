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
		TEST_EMAIL     = "test@example.com"
		TEST_PASSWORD  = "hashedpassword"
		TEST_FULL_NAME = "Test User"
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
				Email:        TEST_EMAIL,
				FullName:     TEST_FULL_NAME,
				PasswordHash: TEST_PASSWORD,
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
				FullName:     TEST_FULL_NAME,
				PasswordHash: TEST_PASSWORD,
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
				FullName:     TEST_FULL_NAME,
				PasswordHash: TEST_PASSWORD,
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
				Email:        TEST_EMAIL,
				FullName:     TEST_FULL_NAME,
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
				Email:        TEST_EMAIL,
				FullName:     "",
				PasswordHash: TEST_PASSWORD,
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
				Email:        TEST_EMAIL,
				FullName:     TEST_FULL_NAME,
				PasswordHash: TEST_PASSWORD,
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
				Email:        TEST_EMAIL,
				FullName:     TEST_FULL_NAME,
				PasswordHash: TEST_PASSWORD,
				CreatedAt:    time.Time{},
				UpdatedAt:    time.Time{},
			},
			wantErr: false,
		},
		{
			name: "missing optional fields",
			user: User{
				ID:           testUUID,
				Email:        TEST_EMAIL,
				FullName:     TEST_FULL_NAME,
				PasswordHash: TEST_PASSWORD,
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

		var unmarshaled map[string]interface{}
		err = json.Unmarshal(jsonData, &unmarshaled)
		assert.NoError(t, err)

		_, exists := unmarshaled["password_hash"]
		assert.False(t, exists)

		assert.Equal(t, user.ID.String(), unmarshaled["id"])
		assert.Equal(t, user.Email, unmarshaled["email"])
		assert.Equal(t, user.FullName, unmarshaled["full_name"])
		assert.Equal(t, user.CreatedAt.Format(time.RFC3339Nano), unmarshaled["created_at"])
		assert.Equal(t, user.UpdatedAt.Format(time.RFC3339Nano), unmarshaled["updated_at"])

	})
}
