// Package user provides user management services.
package user

import (
	"context"

	"github.com/vida-plus/api/internal/auth"
	"github.com/vida-plus/api/models"
)

type UserServiceImpl struct {
	users map[string]*models.User
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{users: make(map[string]*models.User)}
}

func (u *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, ok := u.users[email]
	if !ok {
		return nil, nil // TODO: return error
	}
	return user, nil
}

func (u *UserServiceImpl) CreateUser(ctx context.Context, user *models.User) error {
	u.users[user.Email] = user
	return nil
}

func (u *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.GetByEmail(ctx, email)
}

// UserServiceImpl implements UserStore for AuthService compatibility.
var _ auth.UserStore = (*UserServiceImpl)(nil)
