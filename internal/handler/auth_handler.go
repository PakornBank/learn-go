// Package handler implements HTTP request handlers for authentication endpoints.
package handler

import (
	"context"
	"net/http"

	"github.com/PakornBank/learn-go/internal/model"
	"github.com/PakornBank/learn-go/internal/service"
	"github.com/gin-gonic/gin"
)

// Service defines the authentication operations required by AuthHandler.
type Service interface {
	Register(ctx context.Context, input service.RegisterInput) (*model.User, error)
	Login(ctx context.Context, input service.LoginInput) (string, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

// AuthHandler handles HTTP requests for authentication operations.
type AuthHandler struct {
	service Service
}

// NewAuthHandler creates a new AuthHandler with the provided authentication service.
func NewAuthHandler(s Service) *AuthHandler {
	return &AuthHandler{service: s}
}

// Register handles user registration requests.
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

// Login handles user authentication requests.
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

// GetProfile retrieves the authenticated user's profile.
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
