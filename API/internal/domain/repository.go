package domain

import (
	"context"
)

// Repository defines generic database operations
type Repository interface {
	HealthCheck(ctx context.Context) error
}

// UserRepository defines user-specific database operations
type UserRepository interface {
	Repository
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
}
