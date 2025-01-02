// Package repository provides the data access layer for the application.
// It contains functions and methods for interacting with the database,
// including CRUD operations for various entities. This package abstracts
// the database interactions, providing a clean API for the rest of the
// application to use.
package repository

import (
	"context"

	"github.com/PakornBank/learn-go/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user record into the database.
// It takes a context for managing request-scoped values and cancellation,
// and a pointer to a User model which contains the user data to be inserted.
// It returns an error if the operation fails.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail retrieves a user from the database by their email address.
// It takes a context and an email string as parameters and returns a pointer to a User model and an error.
// If the user is found, it returns the user and a nil error.
// If the user is not found or any other error occurs, it returns nil and the error.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByID retrieves a user from the database by their ID.
// It takes a context and a user ID as parameters and returns a pointer to the User model and an error.
// If the user is found, it returns the user and a nil error.
// If the user is not found or any other error occurs, it returns nil and the error.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
