// Package user provides user management services.
package user

import (
	"context"
	"log/slog"

	"github.com/vida-plus/api/models"
)

// UserServiceImpl implements UserStore interface.
type UserServiceImpl struct {
	users map[string]*models.User
}

func NewUserService() models.UserStore {
	return &UserServiceImpl{users: make(map[string]*models.User)}
}

func (u *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	logger := slog.With(
		slog.String("service", "UserService"),
		slog.String("method", "GetByEmail"),
		slog.String("email", email),
	)

	user, ok := u.users[email]
	if !ok {
		logger.Info("user not found")
		return nil, nil
	}

	logger.Info("user found", slog.String("userID", user.ID))
	return user, nil
}

func (u *UserServiceImpl) Create(ctx context.Context, user *models.User) error {
	logger := slog.With(
		slog.String("service", "UserService"),
		slog.String("method", "Create"),
		slog.String("email", user.Email),
		slog.String("userID", user.ID),
	)

	if _, exists := u.users[user.Email]; exists {
		logger.Error("attempt to create duplicate user")
		return models.NewConflictError("user already exists")
	}

	u.users[user.Email] = user
	logger.Info("user created successfully")
	return nil
}
