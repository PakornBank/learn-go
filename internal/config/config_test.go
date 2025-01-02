package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		want        *Config
		wantErr     bool
		errContains string
	}{
		{
			name:        "default values without JWT secret",
			env:         map[string]string{},
			wantErr:     true,
			errContains: "jwt secret must be set in environment",
		},
		{
			name: "default values with JWT secret",
			env: map[string]string{
				"JWT_SECRET": "test-secret",
			},
			want: &Config{
				DBHost:         "localhost",
				DBUser:         "postgres",
				DBPassword:     "",
				DBName:         "go_auth_db",
				DBPort:         "5432",
				ServerPort:     "8080",
				JWTSecret:      "test-secret",
				TokenExpiryDur: 24 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "custom .env values",
			env: map[string]string{
				"DB_HOST":     "test-db-host",
				"DB_USER":     "test-db-user",
				"DB_PASSWORD": "test-db-password",
				"DB_NAME":     "test-db-name",
				"DB_PORT":     "8081",
				"SERVER_PORT": "5433",
				"JWT_SECRET":  "test-secret",
			},
			want: &Config{
				DBHost:         "test-db-host",
				DBUser:         "test-db-user",
				DBPassword:     "test-db-password",
				DBName:         "test-db-name",
				DBPort:         "8081",
				ServerPort:     "5433",
				JWTSecret:      "test-secret",
				TokenExpiryDur: 24 * time.Hour,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()

			for k, v := range tt.env {
				os.Setenv(k, v)
			}

			got, err := LoadConfig()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		defValue string
		envValue string
		want     string
	}{
		{
			name:     "existing environment variable",
			key:      "TEST_KEY",
			defValue: "default",
			envValue: "custom",
			want:     "custom",
		},
		{
			name:     "non-existing environment variable",
			key:      "TEST_KEY",
			defValue: "default",
			envValue: "",
			want:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			got := getEnv(tt.key, tt.defValue)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDBURL(t *testing.T) {
	cfg := &Config{
		DBHost:     "test-host",
		DBUser:     "test-user",
		DBPassword: "test-password",
		DBName:     "test-name",
		DBPort:     "5432",
	}

	want := "host=test-host user=test-user password=test-password dbname=test-name port=5432 sslmode=disable"
	assert.Equal(t, want, cfg.DBURL())
}
