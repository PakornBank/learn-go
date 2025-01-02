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

func (ms *MockService) Register(ctx context.Context, in service.RegisterInput) (*model.User, error) {
	args := ms.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (ms *MockService) Login(ctx context.Context, in service.LoginInput) (string, error) {
	args := ms.Called(ctx, in)
	return args.Get(0).(string), args.Error(1)
}

func (ms *MockService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	args := ms.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func setupTest(middleware gin.HandlerFunc) (*gin.Engine, *MockService) {
	gin.SetMode(gin.TestMode)

	mockservice := new(MockService)
	handler := NewAuthHandler(mockservice)

	router := gin.New()
	group := router.Group("/api")
	if middleware != nil {
		group.Use(middleware)
	}
	{
		group.POST("/register", handler.Register)
		group.POST("/login", handler.Login)
		group.GET("/profile", handler.GetProfile)
	}

	return router, mockservice
}

func TestNewAuthHandler(t *testing.T) {
	service := new(MockService)
	handler := NewAuthHandler(service)

	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.service)
}

func TestAuthHandler_Register(t *testing.T) {
	user := testutil.NewMockUser()

	tests := []struct {
		name     string
		input    service.RegisterInput
		mock     func(*MockService)
		wantCode int
		errMsg   string
	}{
		{
			name: "successful registration",
			input: service.RegisterInput{
				Email:    user.Email,
				Password: "password",
				FullName: user.FullName,
			},
			mock: func(ms *MockService) {
				ms.On("Register", mock.Anything, mock.MatchedBy(func(in service.RegisterInput) bool {
					return in.Email == user.Email &&
						in.FullName == user.FullName &&
						in.Password == "password"
				})).Return(&user, nil)
			},
			wantCode: http.StatusCreated,
		},
		{
			name: "auth_service error",
			input: service.RegisterInput{
				Email:    user.Email,
				Password: "password",
				FullName: user.FullName,
			},
			mock: func(ms *MockService) {
				ms.On("Register", mock.Anything, mock.MatchedBy(func(in service.RegisterInput) bool {
					return in.Email == user.Email &&
						in.FullName == user.FullName &&
						in.Password == "password"
				})).Return(nil, errors.New("auth_service error"))
			},
			wantCode: http.StatusBadRequest,
			errMsg:   "auth_service error",
		},
		{
			name: "invalid email",
			input: service.RegisterInput{
				Email:    "",
				Password: "password",
				FullName: user.FullName,
			},
			wantCode: http.StatusBadRequest,
			errMsg:   "Error:Field validation for 'Email' failed",
		},
		{
			name: "invalid password",
			input: service.RegisterInput{
				Email:    user.Email,
				Password: "",
				FullName: user.FullName,
			},
			wantCode: http.StatusBadRequest,
			errMsg:   "Error:Field validation for 'Password' failed",
		},
		{
			name: "invalid full name",
			input: service.RegisterInput{
				Email:    user.Email,
				Password: "password",
				FullName: "",
			},
			wantCode: http.StatusBadRequest,
			errMsg:   "Error:Field validation for 'FullName' failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupTest(nil)
			if tt.mock != nil {
				tt.mock(mockService)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)

			assert.NotNil(t, res)

			if tt.wantCode == http.StatusCreated {
				assert.Equal(t, user.ID.String(), res["id"])
				assert.Equal(t, user.FullName, res["full_name"])
				assert.Equal(t, user.Email, res["email"])
				assert.Equal(t, user.CreatedAt.Format(time.RFC3339Nano), res["created_at"])
				assert.Equal(t, user.UpdatedAt.Format(time.RFC3339Nano), res["updated_at"])
				assert.Empty(t, res["password_hash"])
			} else {
				assert.Contains(t, res["error"], tt.errMsg)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	const (
		testToken    = "test-token"
		testEmail    = "test@example.com"
		testPassword = "password"
	)

	tests := []struct {
		name     string
		input    service.LoginInput
		mock     func(*MockService)
		wantCode int
		errMsg   string
	}{
		{
			name: "successful login",
			input: service.LoginInput{
				Email:    testEmail,
				Password: testPassword,
			},
			mock: func(ms *MockService) {
				ms.On("Login", mock.Anything, mock.MatchedBy(func(input service.LoginInput) bool {
					return input.Email == testEmail && input.Password == testPassword
				})).Return(testToken, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "auth_service error",
			input: service.LoginInput{
				Email:    testEmail,
				Password: testPassword,
			},
			mock: func(ms *MockService) {
				ms.On("Login", mock.Anything, mock.MatchedBy(func(input service.LoginInput) bool {
					return input.Email == testEmail && input.Password == testPassword
				})).Return("", errors.New("auth_service error"))
			},
			wantCode: http.StatusBadRequest,
			errMsg:   "auth_service error",
		},
		{
			name: "invalid email",
			input: service.LoginInput{
				Email:    "",
				Password: testPassword,
			},
			wantCode: http.StatusBadRequest,
			errMsg:   "Error:Field validation for 'Email' failed",
		},
		{
			name: "invalid password",
			input: service.LoginInput{
				Email:    testEmail,
				Password: "",
			},
			wantCode: http.StatusBadRequest,
			errMsg:   "Error:Field validation for 'Password' failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupTest(nil)
			if tt.mock != nil {
				tt.mock(mockService)
			}

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var res map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &res)
			assert.NoError(t, err)

			assert.NotNil(t, res)

			if tt.wantCode == http.StatusOK {
				assert.Equal(t, testToken, res["token"])
			} else {
				assert.Contains(t, res["error"], tt.errMsg)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_GetProfile(t *testing.T) {
	user := testutil.NewMockUser()

	tests := []struct {
		name       string
		middleware gin.HandlerFunc
		mock       func(*MockService)
		wantCode   int
		errMsg     string
	}{
		{
			name: "successful profile retrieval",
			middleware: func(c *gin.Context) {
				c.Set("user_id", user.ID.String())
			},
			mock: func(ms *MockService) {
				ms.On("GetUserByID", mock.Anything, user.ID.String()).
					Return(&user, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "auth_service error",
			middleware: func(c *gin.Context) {
				c.Set("user_id", user.ID.String())
			},
			mock: func(ms *MockService) {
				ms.On("GetUserByID", mock.Anything, user.ID.String()).
					Return(nil, errors.New("auth_service error"))
			},
			wantCode: http.StatusNotFound,
			errMsg:   "auth_service error",
		},
		{
			name:     "no user_id in context",
			wantCode: http.StatusUnauthorized,
			errMsg:   "unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockService := setupTest(tt.middleware)
			if tt.mock != nil {
				tt.mock(mockService)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/profile", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			if tt.wantCode == http.StatusOK {
				var res model.User
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)

				assert.Equal(t, user.ID, res.ID)
				assert.Equal(t, user.Email, res.Email)
				assert.Equal(t, user.FullName, res.FullName)
				assert.Equal(t, user.CreatedAt.Format(time.RFC3339Nano), res.CreatedAt.Format(time.RFC3339Nano))
				assert.Equal(t, user.UpdatedAt.Format(time.RFC3339Nano), res.UpdatedAt.Format(time.RFC3339Nano))
				assert.Empty(t, res.PasswordHash)
			} else {
				var res map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &res)
				assert.NoError(t, err)

				assert.Contains(t, res["error"], tt.errMsg)
			}

			mockService.AssertExpectations(t)
		})
	}
}
