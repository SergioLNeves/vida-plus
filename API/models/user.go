// Package models contains domain models for authentication and user management.
package models

import "context"

// User represents a user in the system.
type User struct {
	ID       string
	Email    string
	Password string // hashed
}

// UserStore defines user management methods.
type UserStore interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, user *User) error
}
