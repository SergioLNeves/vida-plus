// Package user provides user management services.
package user

import (
	"context"
	"log/slog"

	"github.com/vida-plus/api/models"
)

// UserServiceImpl implements UserStore interface.
type UserServiceImpl struct {
	repo models.UserRepository
}

func NewUserService(repo models.UserRepository) models.UserStore {
	return &UserServiceImpl{repo: repo}
}

func (u *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	logger := slog.With(
		slog.String("service", "UserService"),
		slog.String("method", "GetByEmail"),
		slog.String("email", email),
	)

	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		logger.Error("failed to get user by email", slog.Any("error", err))
		return nil, err
	}

	if user == nil {
		logger.Info("user not found")
		return nil, nil
	}

	logger.Info("user found successfully", slog.String("userID", user.ID))
	return user, nil
}

func (u *UserServiceImpl) Create(ctx context.Context, user *models.User) error {
	logger := slog.With(
		slog.String("service", "UserService"),
		slog.String("method", "Create"),
		slog.String("email", user.Email),
		slog.String("userID", user.ID),
	)

	if err := u.repo.CreateUser(ctx, user); err != nil {
		logger.Error("failed to create user", slog.Any("error", err))
		return err
	}

	logger.Info("user created successfully")
	return nil
}
