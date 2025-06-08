package pkg

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/vida-plus/api/models"
)

func TestJWTManager_Generate(t *testing.T) {
	t.Run("should generate valid JWT token", func(t *testing.T) {
		// Setup
		secret := "test-secret"
		jwtManager := NewJWTManager(secret)

		user := &models.User{
			ID:    "test-user-id",
			Email: "test@example.com",
		}

		// Execute
		token, err := jwtManager.Generate(user)

		// Verify
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token can be parsed
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		assert.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		// Verify claims
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, user.ID, claims["user_id"])
		assert.Equal(t, user.Email, claims["email"])

		// Verify expiration is in the future
		exp, ok := claims["exp"].(float64)
		assert.True(t, ok)
		assert.True(t, time.Unix(int64(exp), 0).After(time.Now()))
	})

	t.Run("should generate different tokens for different users", func(t *testing.T) {
		// Setup
		secret := "test-secret"
		jwtManager := NewJWTManager(secret)

		user1 := &models.User{
			ID:    "user-1",
			Email: "user1@example.com",
		}

		user2 := &models.User{
			ID:    "user-2",
			Email: "user2@example.com",
		}

		// Execute
		token1, err1 := jwtManager.Generate(user1)
		token2, err2 := jwtManager.Generate(user2)

		// Verify
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestJWTManager_Validate(t *testing.T) {
	t.Run("should validate valid JWT token", func(t *testing.T) {
		// Setup
		secret := "test-secret"
		jwtManager := NewJWTManager(secret)

		user := &models.User{
			ID:    "test-user-id",
			Email: "test@example.com",
		}

		// Generate token
		token, err := jwtManager.Generate(user)
		assert.NoError(t, err)

		// Execute
		claims, err := jwtManager.Validate(token)

		// Verify
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, user.ID, claims.UserID)
		assert.Equal(t, user.Email, claims.Email)
	})

	t.Run("should return error for invalid token", func(t *testing.T) {
		// Setup
		secret := "test-secret"
		jwtManager := NewJWTManager(secret)

		invalidToken := "invalid.jwt.token"

		// Execute
		claims, err := jwtManager.Validate(invalidToken)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("should return error for token with wrong secret", func(t *testing.T) {
		// Setup
		secret1 := "secret-1"
		secret2 := "secret-2"

		jwtManager1 := NewJWTManager(secret1)
		jwtManager2 := NewJWTManager(secret2)

		user := &models.User{
			ID:    "test-user-id",
			Email: "test@example.com",
		}

		// Generate token with secret1
		token, err := jwtManager1.Generate(user)
		assert.NoError(t, err)

		// Try to validate with secret2
		claims, err := jwtManager2.Validate(token)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("should return error for expired token", func(t *testing.T) {
		// Setup
		secret := "test-secret"

		// Create expired token manually
		claims := jwt.MapClaims{
			"user_id": "test-user-id",
			"email":   "test@example.com",
			"exp":     time.Now().Add(-1 * time.Hour).Unix(), // expired 1 hour ago
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		assert.NoError(t, err)

		jwtManager := NewJWTManager(secret)

		// Execute
		validatedClaims, err := jwtManager.Validate(tokenString)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, validatedClaims)
		assert.Contains(t, err.Error(), "invalid token")
	})

	t.Run("should return error for malformed token", func(t *testing.T) {
		// Setup
		secret := "test-secret"
		jwtManager := NewJWTManager(secret)

		malformedToken := "not.a.jwt"

		// Execute
		claims, err := jwtManager.Validate(malformedToken)

		// Verify
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Contains(t, err.Error(), "invalid token")
	})
}
