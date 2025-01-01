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

func generateTestToken(userID string, email string, jwtSecret string, expiry time.Duration) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(expiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(jwtSecret))
	return signedToken
}

func TestAuthMiddleware(t *testing.T) {
	const (
		testSecret = "test-secret"
		testID     = "test-user-id"
		testEmail  = "test@email.com"
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
				return "Bearer " + generateTestToken(testID, testEmail, testSecret, time.Hour)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "expired token",
			generateAuthHeader: func() string {
				return "Bearer " + generateTestToken(testID, testEmail, testSecret, -time.Hour)
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
				return generateTestToken(testID, testEmail, testSecret, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid authorization header format",
		},
		{
			name: "wrong signing method",
			generateAuthHeader: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
					"user_id": testID,
					"email":   testEmail,
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
					"id": testID,
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				signedToken, _ := token.SignedString([]byte(testSecret))
				return "Bearer " + signedToken
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
		{
			name: "missing user_id claim",
			generateAuthHeader: func() string {
				return "Bearer " + generateTestToken("", testEmail, testSecret, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
		{
			name: "missing email claim",
			generateAuthHeader: func() string {
				return "Bearer " + generateTestToken(testID, "", testSecret, time.Hour)
			},
			wantCode:    http.StatusUnauthorized,
			errContains: "invalid token claims",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTest(testSecret)

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
				assert.Equal(t, testID, response["user_id"])
				assert.Equal(t, testEmail, response["email"])
			} else {
				assert.Contains(t, response["error"], tt.errContains)
			}
		})
	}
}
