// Package handler implements HTTP request handlers for authentication endpoints.
package handler

import (
	"context"
	"net/http"

	"github.com/PakornBank/learn-go/internal/model"
	"github.com/PakornBank/learn-go/internal/service"
	"github.com/gin-gonic/gin"
)

// Service defines the methods that an authentication handler must implement.
// It includes methods for user registration, login, and retrieving user information by ID.
type Service interface {
	// Register registers a new user with the given input and returns the created user or an error.
	// ctx: The context for the request.
	// input: The input data required for user registration.
	Register(ctx context.Context, input service.RegisterInput) (*model.User, error)

	// Login authenticates a user with the given input and returns a token or an error.
	// ctx: The context for the request.
	// input: The input data required for user login.
	Login(ctx context.Context, input service.LoginInput) (string, error)

	// GetUserByID retrieves a user by their ID and returns the user or an error.
	// ctx: The context for the request.
	// id: The ID of the user to retrieve.
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

// AuthHandler handles authentication-related HTTP requests.
// It uses a Service to perform the necessary authentication operations.
type AuthHandler struct {
	service Service
}

// NewAuthHandler creates a new instance of AuthHandler with the provided service.
// It returns a pointer to the newly created AuthHandler.
//
// Parameters:
//   - s: The service that the AuthHandler will use.
//
// Returns:
//   - A pointer to the newly created AuthHandler.
func NewAuthHandler(s Service) *AuthHandler {
	return &AuthHandler{service: s}
}

// Register handles the user registration process.
// It binds the JSON input to the RegisterInput struct and calls the service's Register method.
// If the input is invalid or the registration fails, it responds with a 400 status code and an error message.
// On successful registration, it responds with a 201 status code and the created user.
func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login handles the user login process.
// It expects a JSON payload with login credentials, binds it to a LoginInput struct,
// and attempts to authenticate the user using the AuthService.
// If successful, it returns a JSON response with an authentication token.
// If there is an error during binding or authentication, it returns a JSON response with the error message.
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.service.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// GetProfile handles the request to retrieve the profile of the authenticated user.
// It expects the user ID to be stored in the context with the key "user_id".
// If the user ID is not found in the context, it responds with an unauthorized status.
// If the user ID is found, it attempts to retrieve the user profile from the service.
// If the user profile is not found, it responds with a not found status.
// If the user profile is successfully retrieved, it responds with the user profile in JSON format.
func (h *AuthHandler) GetProfile(c *gin.Context) {
	id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), id.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
