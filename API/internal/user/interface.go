// Package user provides user management services.
package user

import (
	"context"

	"github.com/vida-plus/api/models"
)

// UserService defines user management methods.
type UserService interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
