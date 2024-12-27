package service

import (
	"context"
	"testing"
	"time"

	"github.com/PakornBank/learn-go/internal/config"
	"github.com/PakornBank/learn-go/internal/model"
	"github.com/PakornBank/learn-go/internal/testutil"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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

func setupTest(t *testing.T) (*AuthService, *MockRepository) {
	mockRepo := new(MockRepository)
	config := &config.Config{
		JWTSecret:      "test-secret",
		TokenExipryDur: time.Hour * 24,
	}
	service := NewAuthService(mockRepo, config)
	return service, mockRepo
}

func TestRegister(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name        string
		input       RegisterInput
		mockFn      func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful registration",
			input: RegisterInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(repo *MockRepository) {
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(nil, gorm.ErrRecordNotFound)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "email already exists",
			input: RegisterInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(repo *MockRepository) {
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(&mockUser, nil)
			},
			wantErr:     true,
			errContains: "email already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTest(t)
			tt.mockFn(mockRepo)
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
	mockUser := testutil.NewMockUser()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	tests := []struct {
		name        string
		input       LoginInput
		mockFn      func(*MockRepository)
		wantErr     bool
		errContains string
	}{
		{
			name: "successful login",
			input: LoginInput{
				Email:    mockUser.Email,
				Password: "password",
			},
			mockFn: func(repo *MockRepository) {
				mockUser.PasswordHash = string(hashedPassword)
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(&mockUser, nil)
			},
			wantErr: false,
		},
		{
			name: "invalid credentials",
			input: LoginInput{
				Email:    mockUser.Email,
				Password: "wrongpassword",
			},
			mockFn: func(repo *MockRepository) {
				mockUser.PasswordHash = string(hashedPassword)
				repo.On("FindByEmail", mock.Anything, mockUser.Email).Return(&mockUser, nil)
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
		{
			name: "user not found",
			input: LoginInput{
				Email:    "nonexistent@example.com",
				Password: "password",
			},
			mockFn: func(repo *MockRepository) {
				repo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr:     true,
			errContains: "invalid credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTest(t)
			tt.mockFn(mockRepo)
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
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name    string
		id      string
		mockFn  func(*MockRepository)
		want    *model.User
		wantErr bool
		errType error
	}{
		{
			name: "user found",
			id:   mockUser.ID.String(),
			mockFn: func(repo *MockRepository) {
				repo.On("FindById", mock.Anything, mockUser.ID.String()).Return(&mockUser, nil)
			},
			want:    &mockUser,
			wantErr: false,
		},
		{
			name: "user not found",
			id:   mockUser.ID.String(),
			mockFn: func(repo *MockRepository) {
				repo.On("FindById", mock.Anything, mockUser.ID.String()).Return(nil, gorm.ErrRecordNotFound)
			},
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTest(t)
			tt.mockFn(mockRepo)
			got, err := service.GetUserById(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGenerateToken(t *testing.T) {
	service, _ := setupTest(t)
	mockUser := testutil.NewMockUser()

	token, err := service.generateToken(&mockUser)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, mockUser.ID.String(), claims["user_id"])
	assert.Equal(t, mockUser.Email, claims["email"])
}
