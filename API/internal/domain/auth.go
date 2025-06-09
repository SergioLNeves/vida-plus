// Package models contains domain models for authentication and authorization.
package domain

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrUnauthorized = errors.New("unauthorized access")

// AuthClaims represents JWT claims for authentication.
type AuthClaims struct {
	UserID   string   `json:"user_id"`
	Email    string   `json:"email"`
	UserType UserType `json:"user_type"`
	jwt.RegisteredClaims
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
	RegisterWithProfile(ctx context.Context, req RegisterRequest) (*User, error)
	Login(ctx context.Context, email, password string) (token string, err error)
	GenerateRefreshToken(ctx context.Context, user *User) (string, error)
	ValidateRefreshToken(ctx context.Context, token string) (*User, error)
}

// JWTManager defines methods for JWT token management.
type JWTManager interface {
	Generate(user *User) (string, error)
	Validate(token string) (*AuthClaims, error)
	GenerateRefreshToken(user *User) (string, error)
	ValidateRefreshToken(token string) (*AuthClaims, error)
}
