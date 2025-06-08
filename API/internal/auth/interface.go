// Package auth provides authentication and authorization services.
package auth

import (
	"context"

	"github.com/vida-plus/api/models"
)

// AuthService defines authentication methods.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (token string, err error)
}
