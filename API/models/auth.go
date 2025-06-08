// Package models contains domain models for authentication and authorization.
package models

import (
	"context"
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized access")

// AuthClaims represents JWT claims for authentication.
type AuthClaims struct {
	UserID string
	Email  string
}

// GetAuthClaims extracts JWT claims from echo.Context
func GetAuthClaims(claims interface{}) (*AuthClaims, error) {
	if claims == nil {
		return nil, ErrUnauthorized
	}

	authClaims, ok := claims.(*AuthClaims)
	if !ok {
		return nil, ErrUnauthorized
	}

	return authClaims, nil
}

// AuthService defines authentication methods.
type AuthService interface {
	Register(ctx context.Context, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (token string, err error)
}

// JWTManager defines methods for JWT token management.
type JWTManager interface {
	Generate(user *User) (string, error)
	Validate(token string) (*AuthClaims, error)
}
