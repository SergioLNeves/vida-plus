// Package auth provides authentication and authorization services.
package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/vida-plus/api/models"
)

type AuthServiceImpl struct {
	userStore UserStore
	jwt       JWTManager
}

// UserStore abstracts user persistence.
type UserStore interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

// JWTManager abstracts JWT operations.
type JWTManager interface {
	Generate(user *models.User) (string, error)
	Verify(token string) (*models.AuthClaims, error)
}

func NewAuthService(userStore UserStore, jwt JWTManager) AuthService {
	return &AuthServiceImpl{userStore: userStore, jwt: jwt}
}

func (a *AuthServiceImpl) Register(ctx context.Context, email, password string) (*models.User, error) {
	// TODO: hash password
	user := &models.User{Email: email, Password: password}
	if err := a.userStore.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := a.userStore.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("invalid credentials")
	}
	return a.jwt.Generate(user)
}
