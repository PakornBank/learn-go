package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PakornBank/learn-go/internal/model"
	"github.com/PakornBank/learn-go/internal/service"
	"github.com/PakornBank/learn-go/internal/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Register(ctx context.Context, input service.RegisterInput) (*model.User, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockService) Login(ctx context.Context, input service.LoginInput) (string, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockService) GetUserById(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func setupTest(middleware gin.HandlerFunc) (*gin.Engine, *MockService) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	handler := NewAuthHandler(mockService)

	router := gin.New()
	r := router.Group("/api")
	if middleware != nil {
		r.Use(middleware)
	}
	{
		r.POST("/register", handler.Register)
		r.POST("/login", handler.Login)
		r.GET("/profile", handler.GetProfile)
	}

	return router, mockService
}

func TestNewAuthHandler(t *testing.T) {
	mockService := new(MockService)
	handler := NewAuthHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.authService)
}

func TestAuthHandler_Register(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name        string
		input       service.RegisterInput
		mockFn      func(*MockService)
		wantCode    int
		errContains string
	}{
		{
			name: "successful registration",
			input: service.RegisterInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(ms *MockService) {
				ms.On("Register", mock.Anything, mock.MatchedBy(func(input service.RegisterInput) bool {
					return input.Email == mockUser.Email &&
						input.FullName == mockUser.FullName &&
						input.Password == "password"
				})).Return(&mockUser, nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "auth_service error",
			input: service.RegisterInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: mockUser.FullName,
			},
			mockFn: func(ms *MockService) {
				ms.On("Register", mock.Anything, mock.MatchedBy(func(input service.RegisterInput) bool {
					return input.Email == mockUser.Email &&
						input.FullName == mockUser.FullName &&
						input.Password == "password"
				})).Return(nil, errors.New("auth_service error"))
			},
			wantCode:    http.StatusBadRequest,
			errContains: "auth_service error",
		},
		{
			name: "invalid email",
			input: service.RegisterInput{
				Email:    "",
				Password: "password",
				FullName: mockUser.FullName,
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Email' failed",
		},
		{
			name: "invalid password",
			input: service.RegisterInput{
				Email:    mockUser.Email,
				Password: "",
				FullName: mockUser.FullName,
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Password' failed",
		},
		{
			name: "invalid full name",
			input: service.RegisterInput{
				Email:    mockUser.Email,
				Password: "password",
				FullName: "",
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'FullName' failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupTest(nil)
			if tt.mockFn != nil {
				tt.mockFn(mockService)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.NotNil(t, response)

			if tt.wantCode == http.StatusCreated {
				assert.Equal(t, mockUser.ID.String(), response["id"])
				assert.Equal(t, mockUser.FullName, response["full_name"])
				assert.Equal(t, mockUser.Email, response["email"])
				assert.Equal(t, mockUser.CreatedAt.Format(time.RFC3339Nano), response["created_at"])
				assert.Equal(t, mockUser.UpdatedAt.Format(time.RFC3339Nano), response["updated_at"])
				assert.Empty(t, response["password_hash"])
			} else {
				assert.Contains(t, response["error"], tt.errContains)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	const (
		TEST_TOKEN    = "test-token"
		TEST_EMAIL    = "test@example.com"
		TEST_PASSWORD = "password"
	)

	tests := []struct {
		name        string
		input       service.LoginInput
		mockFn      func(*MockService)
		wantCode    int
		errContains string
	}{
		{
			name: "successful login",
			input: service.LoginInput{
				Email:    TEST_EMAIL,
				Password: TEST_PASSWORD,
			},
			mockFn: func(ms *MockService) {
				ms.On("Login", mock.Anything, mock.MatchedBy(func(input service.LoginInput) bool {
					return input.Email == TEST_EMAIL && input.Password == TEST_PASSWORD
				})).Return(TEST_TOKEN, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "auth_service error",
			input: service.LoginInput{
				Email:    TEST_EMAIL,
				Password: TEST_PASSWORD,
			},
			mockFn: func(ms *MockService) {
				ms.On("Login", mock.Anything, mock.MatchedBy(func(input service.LoginInput) bool {
					return input.Email == TEST_EMAIL && input.Password == TEST_PASSWORD
				})).Return("", errors.New("auth_service error"))
			},
			wantCode:    http.StatusBadRequest,
			errContains: "auth_service error",
		},
		{
			name: "invalid email",
			input: service.LoginInput{
				Email:    "",
				Password: TEST_PASSWORD,
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Email' failed",
		},
		{
			name: "invalid password",
			input: service.LoginInput{
				Email:    TEST_EMAIL,
				Password: "",
			},
			wantCode:    http.StatusBadRequest,
			errContains: "Error:Field validation for 'Password' failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupTest(nil)
			if tt.mockFn != nil {
				tt.mockFn(mockService)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.NotNil(t, response)

			if tt.wantCode == http.StatusOK {
				assert.Equal(t, TEST_TOKEN, response["token"])
			} else {
				assert.Contains(t, response["error"], tt.errContains)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	mockUser := testutil.NewMockUser()

	tests := []struct {
		name        string
		middleware  gin.HandlerFunc
		mockFn      func(*MockService)
		wantCode    int
		errContains string
	}{
		{
			name: "successful profile retrieval",
			middleware: func(c *gin.Context) {
				c.Set("user_id", mockUser.ID.String())
			},
			mockFn: func(ms *MockService) {
				ms.On("GetUserById", mock.Anything, mockUser.ID.String()).
					Return(&mockUser, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "auth_service error",
			middleware: func(c *gin.Context) {
				c.Set("user_id", mockUser.ID.String())
			},
			mockFn: func(ms *MockService) {
				ms.On("GetUserById", mock.Anything, mockUser.ID.String()).
					Return(nil, errors.New("auth_service error"))
			},
			wantCode:    http.StatusNotFound,
			errContains: "auth_service error",
		},
		{
			name:        "no user_id in context",
			wantCode:    http.StatusUnauthorized,
			errContains: "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupTest(tt.middleware)
			if tt.mockFn != nil {
				tt.mockFn(mockService)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode == http.StatusOK {
				var response model.User
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Equal(t, mockUser.ID, response.ID)
				assert.Equal(t, mockUser.Email, response.Email)
				assert.Equal(t, mockUser.FullName, response.FullName)
				assert.Equal(t, mockUser.CreatedAt.Format(time.RFC3339Nano), response.CreatedAt.Format(time.RFC3339Nano))
				assert.Equal(t, mockUser.UpdatedAt.Format(time.RFC3339Nano), response.UpdatedAt.Format(time.RFC3339Nano))
				assert.Empty(t, response.PasswordHash)
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				assert.Contains(t, response["error"], tt.errContains)
			}

			mockService.AssertExpectations(t)
		})
	}
}
