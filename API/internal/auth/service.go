// Package auth provides authentication and authorization services.
package auth

import (
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/vida-plus/api/models"
	"github.com/vida-plus/api/pkg"
)

// AuthServiceImpl implements AuthService interface.
type AuthServiceImpl struct {
	userStore models.UserStore
	jwt       models.JWTManager
}

func NewAuthService(userStore models.UserStore, jwt models.JWTManager) models.AuthService {
	return &AuthServiceImpl{userStore: userStore, jwt: jwt}
}

func (a *AuthServiceImpl) Register(ctx context.Context, email, password string) (*models.User, error) {
	logger := slog.With(
		slog.String("service", "AuthService"),
		slog.String("method", "Register"),
		slog.String("email", email),
	)

	// Check if user already exists
	existingUser, err := a.userStore.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("error checking existing user", slog.Any("error", err))
		return nil, models.NewInternalError("error checking user existence")
	}
	if existingUser != nil {
		logger.Info("attempt to register existing user")
		return nil, models.NewConflictError("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("error hashing password", slog.Any("error", err))
		return nil, models.NewInternalError("error processing password")
	}

	user := &models.User{
		ID:       pkg.GenerateID(),
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := a.userStore.Create(ctx, user); err != nil {
		logger.Error("error creating user", slog.Any("error", err))
		return nil, models.NewInternalError("error creating user")
	}

	logger.Info("user registered successfully", slog.String("userID", user.ID))
	return user, nil
}

func (a *AuthServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	logger := slog.With(
		slog.String("service", "AuthService"),
		slog.String("method", "Login"),
		slog.String("email", email),
	)

	user, err := a.userStore.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("error fetching user", slog.Any("error", err))
		return "", models.NewInternalError("error processing login")
	}
	if user == nil {
		logger.Info("login attempt with non-existent user")
		return "", models.NewUnauthorizedError("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Info("login attempt with invalid password")
		return "", models.NewUnauthorizedError("invalid credentials")
	}

	token, err := a.jwt.Generate(user)
	if err != nil {
		logger.Error("error generating token", slog.Any("error", err))
		return "", models.NewInternalError("error generating authentication token")
	}

	logger.Info("user logged in successfully", slog.String("userID", user.ID))
	return token, nil
}
