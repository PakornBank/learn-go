// Package middleware provides various middleware functions for handling
// authentication and authorization in HTTP requests. These functions
// ensure that requests are properly authenticated and authorized before
// reaching the application logic, enhancing the security and integrity
// of the application.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware is a middleware function for the Gin framework that handles
// JWT authentication. It expects a JWT token in the "Authorization" header
// in the format "Bearer <token>". The token is validated using the provided
// jwtSecret. If the token is valid, the user ID and email from the token
// claims are set in the Gin context.
//
// Parameters:
//   - jwtSecret: The secret key used to validate the JWT token.
//
// Returns:
//   - gin.HandlerFunc: A Gin middleware handler function.
//
// The middleware performs the following checks:
//  1. Ensures the "Authorization" header is present.
//  2. Ensures the "Authorization" header is in the format "Bearer <token>".
//  3. Parses and validates the JWT token using the provided secret.
//  4. Extracts the "user_id" and "email" claims from the token and sets them
//     in the Gin context.
//
// If any of these checks fail, the middleware responds with a 401 Unauthorized
// status and an appropriate error message, and aborts the request.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		userID, hasUserID := claims["user_id"]
		email, hasEmail := claims["email"]
		if !hasUserID || userID == "" || !hasEmail || email == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Set("email", email)
		c.Next()
	}
}
