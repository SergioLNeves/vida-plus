package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/vida-plus/api/mocks"
	"github.com/vida-plus/api/models"
)

func TestAuthService_Register(t *testing.T) {
	t.Run("should register user successfully when user does not exist", func(t *testing.T) {
		// Setup mocks
		mockUserStore := mocks.NewUserStore(t)
		mockJWTManager := mocks.NewJWTManager(t)

		// Setup service
		authService := NewAuthService(mockUserStore, mockJWTManager)

		// Setup test data
		ctx := context.Background()
		email := "test@example.com"
		password := "password123"

		// Configure expectations
		mockUserStore.
			EXPECT().
			GetByEmail(ctx, email).
			Return(nil, nil)

		mockUserStore.
			EXPECT().
			Create(ctx, mock.MatchedBy(func(user *models.User) bool {
				return user.Email == email && user.Password != password && user.ID != ""
			})).
			Return(nil)

		// Execute test
		user, err := authService.Register(ctx, email, password)

		// Verify results
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.NotEqual(t, password, user.Password) // password should be hashed
		assert.NotEmpty(t, user.ID)                 // ID should be generated

		// Verify password is properly hashed
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		assert.NoError(t, err)
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		// Setup mocks
		mockUserStore := mocks.NewUserStore(t)
		mockJWTManager := mocks.NewJWTManager(t)

		// Setup service
		authService := NewAuthService(mockUserStore, mockJWTManager)

		// Setup test data
		ctx := context.Background()
		email := "test@example.com"
		password := "password123"
		existingUser := &models.User{Email: email, Password: "hashedpassword"}

		// Configure expectations
		mockUserStore.
			EXPECT().
			GetByEmail(ctx, email).
			Return(existingUser, nil)

		// Execute test
		user, err := authService.Register(ctx, email, password)

		// Verify results
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user already exists")
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("should login successfully with valid credentials", func(t *testing.T) {
		// Setup mocks
		mockUserStore := mocks.NewUserStore(t)
		mockJWTManager := mocks.NewJWTManager(t)

		// Setup service
		authService := NewAuthService(mockUserStore, mockJWTManager)

		// Setup test data
		ctx := context.Background()
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &models.User{
			ID:       "user123",
			Email:    email,
			Password: string(hashedPassword),
		}

		expectedToken := "jwt.token.here"

		// Configure expectations
		mockUserStore.
			EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		mockJWTManager.
			EXPECT().
			Generate(user).
			Return(expectedToken, nil)

		// Execute test
		token, err := authService.Login(ctx, email, password)

		// Verify results
		assert.NoError(t, err)
		assert.Equal(t, expectedToken, token)
	})

	t.Run("should return error when user does not exist", func(t *testing.T) {
		// Setup mocks
		mockUserStore := mocks.NewUserStore(t)
		mockJWTManager := mocks.NewJWTManager(t)

		// Setup service
		authService := NewAuthService(mockUserStore, mockJWTManager)

		// Setup test data
		ctx := context.Background()
		email := "test@example.com"
		password := "password123"

		// Configure expectations
		mockUserStore.
			EXPECT().
			GetByEmail(ctx, email).
			Return(nil, nil)

		// Execute test
		token, err := authService.Login(ctx, email, password)

		// Verify results
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "invalid credentials")
	})

	t.Run("should return error with invalid password", func(t *testing.T) {
		// Setup mocks
		mockUserStore := mocks.NewUserStore(t)
		mockJWTManager := mocks.NewJWTManager(t)

		// Setup service
		authService := NewAuthService(mockUserStore, mockJWTManager)

		// Setup test data
		ctx := context.Background()
		email := "test@example.com"
		password := "wrongpassword"
		correctPassword := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

		user := &models.User{
			ID:       "user123",
			Email:    email,
			Password: string(hashedPassword),
		}

		// Configure expectations
		mockUserStore.
			EXPECT().
			GetByEmail(ctx, email).
			Return(user, nil)

		// Execute test
		token, err := authService.Login(ctx, email, password)

		// Verify results
		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Contains(t, err.Error(), "invalid credentials")
	})
}
