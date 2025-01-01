package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthHandler struct {
	mock.Mock
}

func (m *MockAuthHandler) Register(c *gin.Context) {
	m.Called(c)
}

func (m *MockAuthHandler) Login(c *gin.Context) {
	m.Called(c)
}

func (m *MockAuthHandler) GetProfile(c *gin.Context) {
	m.Called(c)
}

func setupTest() (*gin.Engine, *MockAuthHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockHandler := new(MockAuthHandler)
	setupAuthRoutes(router, mockHandler)
	return router, mockHandler
}

func TestSetupAuthRoutes(t *testing.T) {
	tests := []struct {
		name     string
		method   string
		path     string
		mockFn   func(*MockAuthHandler)
		wantCode int
	}{
		{
			name:   "register route",
			method: http.MethodPost,
			path:   "/api/register",
			mockFn: func(m *MockAuthHandler) {
				m.On("Register", mock.Anything).Return()
			},
			wantCode: http.StatusOK,
		},
		{
			name:   "login route",
			method: http.MethodPost,
			path:   "/api/login",
			mockFn: func(mah *MockAuthHandler) {
				mah.On("Login", mock.Anything).Return()
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "undefined route",
			method:   http.MethodGet,
			path:     "/api/undefined",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockHandler := setupTest()
			if tt.mockFn != nil {
				tt.mockFn(mockHandler)
			}

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
			mockHandler.AssertExpectations(t)
		})
	}
}
