package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func setupTest(jwtSecret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware(jwtSecret))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id": c.MustGet("user_id"),
			"email":   c.MustGet("email"),
		})
	})
	return r
}

func generateTestToken(userId string, email string, jwtSecret string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(jwtSecret))
	return signedToken
}

func TestAuthMiddleware(t *testing.T) {
	const (
		TEST_SECRET = "test-secret"
		TEST_ID     = "test-user-id"
		TEST_EMAIL  = "test@email.com"
	)

	tests := []struct {
		name               string
		generateAuthHeader func() string
		wantCode           int
		errContains        string
	}{
		{
			name: "valid token",
			generateAuthHeader: func() string {
				return "Bearer " + generateTestToken(TEST_ID, TEST_EMAIL, TEST_SECRET, time.Hour)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "expired token",
			generateAuthHeader: func() string {
				return "Bearer " + generateTestToken(TEST_ID, TEST_EMAIL, TEST_SECRET, -time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token",
		},
		{
			name: "invalid token",
			generateAuthHeader: func() string {
				return "Bearer invalid-token"
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token",
		},
		{
			name: "empty authorization header",
			generateAuthHeader: func() string {
				return ""
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "authorization header required",
		},
		{
			name: "missing Bearer prifix",
			generateAuthHeader: func() string {
				return generateTestToken(TEST_ID, TEST_EMAIL, TEST_SECRET, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid authorization header format",
		},
		{
			name: "wrong signing method",
			generateAuthHeader: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
					"user_id": TEST_ID,
					"email":   TEST_EMAIL,
					"exp":     time.Now().Add(time.Hour).Unix(),
				})
				signedToken, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
				return "Bearer " + signedToken
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token",
		},
		{
			name: "invalid token claims",
			generateAuthHeader: func() string {
				claims := jwt.MapClaims{
					"id": TEST_ID,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte(TEST_SECRET))
				return "Bearer " + signedToken
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
		{
			name: "missing user_id claim",
			generateAuthHeader: func() string {
				return "Bearer " + generateTestToken("", TEST_EMAIL, TEST_SECRET, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
		{
			name: "missing email claim",
			generateAuthHeader: func() string {
				return "Bearer " + generateTestToken(TEST_ID, "", TEST_SECRET, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTest(TEST_SECRET)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if authHeader := tt.generateAuthHeader(); authHeader != "" {
				req.Header.Set("Authorization", authHeader)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.wantCode == http.StatusOK {
				assert.Equal(t, TEST_ID, response["user_id"])
				assert.Equal(t, TEST_EMAIL, response["email"])
			} else {
				assert.Contains(t, response["error"], tt.errContains)
			}
		})
	}
}
