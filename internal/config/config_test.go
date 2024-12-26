package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		wantConfig  *Config
		wantErr     bool
		errContains string
	}{
		{
			name:        "default values without JWT secret",
			envVars:     map[string]string{},
			wantErr:     true,
			errContains: "JWT_SECRET must be set in .env",
		},
		{
			name: "default values with JWT secret",
			envVars: map[string]string{
				"JWT_SECRET": "test-secret",
			},
			wantConfig: &Config{
				DBHost:         "localhost",
				DBUser:         "postgres",
				DBPassword:     "",
				DBName:         "go_auth_db",
				DBPort:         "5432",
				ServerPort:     "8080",
				JWTSecret:      "test-secret",
				TokenExipryDur: 24 * time.Hour,
			},
			wantErr: false,
		},
		{
			name: "custom .env values",
			envVars: map[string]string{
				"DB_HOST":     "test-db-host",
				"DB_USER":     "test-db-user",
				"DB_PASSWORD": "test-db-password",
				"DB_NAME":     "test-db-name",
				"DB_PORT":     "8081",
				"SERVER_PORT": "5433",
				"JWT_SECRET":  "test-secret",
			},
			wantConfig: &Config{
				DBHost:         "test-db-host",
				DBUser:         "test-db-user",
				DBPassword:     "test-db-password",
				DBName:         "test-db-name",
				DBPort:         "8081",
				ServerPort:     "5433",
				JWTSecret:      "test-secret",
				TokenExipryDur: 24 * time.Hour,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			got, err := LoadConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil || err.Error() != tt.errContains {
					t.Errorf("LoadConfig() error = %v, want error containing %v", err, tt.errContains)
				}
				return
			}

			if !(got.DBHost == tt.wantConfig.DBHost &&
				got.DBUser == tt.wantConfig.DBUser &&
				got.DBPassword == tt.wantConfig.DBPassword &&
				got.DBName == tt.wantConfig.DBName &&
				got.DBPort == tt.wantConfig.DBPort &&
				got.ServerPort == tt.wantConfig.ServerPort &&
				got.JWTSecret == tt.wantConfig.JWTSecret &&
				got.TokenExipryDur == tt.wantConfig.TokenExipryDur) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.wantConfig)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "existing environment variable",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "non-existing environment variable",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
			}

			if got := getEnv(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDBURL(t *testing.T) {
	config := &Config{
		DBHost:     "test-host",
		DBUser:     "test-user",
		DBPassword: "test-password",
		DBName:     "test-name",
		DBPort:     "5432",
	}

	want := "host=test-host user=test-user password=test-password dbname=test-name port=5432 sslmode=disable"
	if got := config.GetDBURL(); got != want {
		t.Errorf("GetDBURL() = %v, want %v", got, want)
	}
}
