// Package service provides business logic services.
package service

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vida-plus/api/internal/domain"
	"github.com/vida-plus/api/pkg"
)

// AuthServiceImpl implements AuthService interface.
type AuthServiceImpl struct {
	userStore domain.UserStore
	jwt       domain.JWTManager
}

func NewAuthService(userStore domain.UserStore, jwt domain.JWTManager) domain.AuthService {
	return &AuthServiceImpl{userStore: userStore, jwt: jwt}
}

func (a *AuthServiceImpl) Register(ctx context.Context, email, password string) (*domain.User, error) {
	logger := slog.With(
		slog.String("service", "AuthService"),
		slog.String("method", "Register"),
		slog.String("email", email),
	)

	// Check if user already exists
	existingUser, err := a.userStore.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("error checking existing user", slog.Any("error", err))
		return nil, domain.NewInternalError("error checking user existence")
	}
	if existingUser != nil {
		logger.Info("attempt to register existing user")
		return nil, domain.NewConflictError("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("error hashing password", slog.Any("error", err))
		return nil, domain.NewInternalError("error processing password")
	}

	user := &domain.User{
		ID:        pkg.GenerateID(),
		Email:     email,
		Password:  string(hashedPassword),
		Type:      domain.UserTypePatient, // Default to patient
		Status:    domain.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.userStore.Create(ctx, user); err != nil {
		logger.Error("error creating user", slog.Any("error", err))
		return nil, domain.NewInternalError("error creating user")
	}

	logger.Info("user registered successfully", slog.String("userID", user.ID))
	return user, nil
}

func (a *AuthServiceImpl) RegisterWithProfile(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	logger := slog.With(
		slog.String("service", "AuthService"),
		slog.String("method", "RegisterWithProfile"),
		slog.String("email", req.Email),
		slog.String("type", string(req.Type)),
	)

	// Check if user already exists
	existingUser, err := a.userStore.GetByEmail(ctx, req.Email)
	if err != nil {
		logger.Error("error checking existing user", slog.Any("error", err))
		return nil, domain.NewInternalError("error checking user existence")
	}
	if existingUser != nil {
		logger.Info("attempt to register existing user")
		return nil, domain.NewConflictError("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("error hashing password", slog.Any("error", err))
		return nil, domain.NewInternalError("error processing password")
	}

	user := &domain.User{
		ID:        pkg.GenerateID(),
		Email:     req.Email,
		Password:  string(hashedPassword),
		Type:      req.Type,
		Status:    domain.UserStatusActive,
		Profile:   req.Profile,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.userStore.Create(ctx, user); err != nil {
		logger.Error("error creating user", slog.Any("error", err))
		return nil, domain.NewInternalError("error creating user")
	}

	logger.Info("user registered successfully", slog.String("userID", user.ID), slog.String("type", string(user.Type)))
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
		return "", domain.NewInternalError("error processing login")
	}
	if user == nil {
		logger.Info("login attempt with non-existent user")
		return "", domain.NewUnauthorizedError("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Info("login attempt with invalid password")
		return "", domain.NewUnauthorizedError("invalid credentials")
	}

	token, err := a.jwt.Generate(user)
	if err != nil {
		logger.Error("error generating token", slog.Any("error", err))
		return "", domain.NewInternalError("error generating authentication token")
	}

	logger.Info("user logged in successfully", slog.String("userID", user.ID))
	return token, nil
}

func (a *AuthServiceImpl) GenerateRefreshToken(ctx context.Context, user *domain.User) (string, error) {
	logger := slog.With(
		slog.String("service", "AuthService"),
		slog.String("method", "GenerateRefreshToken"),
		slog.String("userID", user.GetID()),
	)

	refreshToken, err := a.jwt.GenerateRefreshToken(user)
	if err != nil {
		logger.Error("error generating refresh token", slog.Any("error", err))
		return "", domain.NewInternalError("error generating refresh token")
	}

	return refreshToken, nil
}

func (a *AuthServiceImpl) ValidateRefreshToken(ctx context.Context, token string) (*domain.User, error) {
	logger := slog.With(
		slog.String("service", "AuthService"),
		slog.String("method", "ValidateRefreshToken"),
	)

	claims, err := a.jwt.ValidateRefreshToken(token)
	if err != nil {
		logger.Error("error validating refresh token", slog.Any("error", err))
		return nil, domain.NewUnauthorizedError("invalid refresh token")
	}

	user, err := a.userStore.GetByID(ctx, claims.UserID)
	if err != nil {
		logger.Error("error fetching user by ID", slog.Any("error", err))
		return nil, domain.NewInternalError("error fetching user")
	}

	return user, nil
}
