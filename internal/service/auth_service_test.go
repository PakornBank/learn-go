package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) FindById(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func setupTestService() (*AuthService, *MockRepository) {
	mockRepo := new(MockRepository)
	config := &config.Config{
		JWTSecret:      "test-secret",
		TokenExipryDur: time.Hour * 24,
	}
	service := NewAuthService(mockRepo, config)
	return service, mockRepo
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name        string
		input       RegisterInput
		setupMock   func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful registration",
			input: RegisterInput{
				Email:    "test@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			setupMock: func(repo *MockRepository) {
				repo.On("FindByEmail", ctx, "test@example.com").Return(nil, nil)
				repo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			input: RegisterInput{
				Email:    "existing@example.com",
				Password: "password123",
				FullName: "Test User",
			},
			setupMock: func(repo *MockRepository) {
				existingUser := &model.User{Email: "existing@example.com"}
				repo.On("FindByEmail", ctx, "existing@example.com").Return(existingUser, nil)
			},
			wantErr:     true,
			errContains: "email already registered",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTestService()
			tt.setupMock(mockRepo)

			user, err := service.Register(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errContains, err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.input.Email, user.Email)
				assert.Equal(t, tt.input.FullName, user.FullName)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		name        string
		input       LoginInput
		setupMock   func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful login",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(repo *MockRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					ID:           uuid.New(),
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
				}
				repo.On("FindByEmail", ctx, "test@example.com").Return(user, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid credentials - wrong password",
			input: LoginInput{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			setupMock: func(repo *MockRepository) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				user := &model.User{
					Email:        "test@example.com",
					PasswordHash: string(hashedPassword),
				}
				repo.On("FindByEmail", ctx, "test@example.com").Return(user, nil)
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "user not found",
			input: LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			setupMock: func(repo *MockRepository) {
				repo.On("FindByEmail", ctx, "nonexistent@example.com").Return(nil, errors.New("user not found"))
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTestService()
			tt.setupMock(mockRepo)

			token, err := service.Login(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errContains, err.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserById(t *testing.T) {
	testCases := []struct {
		name      string
		userId    string
		setupMock func(*MockRepository)
		wantErr   bool
	}{
		{
			name:   "successful user retrieval",
			userId: "123e4567-e89b-12d3-a456-426614174000",
			setupMock: func(repo *MockRepository) {
				user := &model.User{
					ID:       uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Email:    "test@example.com",
					FullName: "Test User",
				}
				repo.On("FindById", mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(user, nil)
			},
			wantErr: false,
		},
		{
			name:   "user not found",
			userId: "123e4567-e89b-12d3-a456-426614174000",
			setupMock: func(repo *MockRepository) {
				repo.On("FindById", mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(nil, errors.New("user not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTestService()
			tt.setupMock(mockRepo)

			user, err := service.GetUserById(context.Background(), tt.userId)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userId, user.ID.String())
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	service, _ := setupTestService()
	user := &model.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	token, err := service.generateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, user.ID.String(), claims["user_id"])
	assert.Equal(t, user.Email, claims["email"])
}
