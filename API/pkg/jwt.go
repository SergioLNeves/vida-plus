// Package jwt provides JWT utilities for authentication.
package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vida-plus/api/models"
)

// JWTManagerImpl implements JWTManager interface.
type JWTManagerImpl struct {
	secret string
}

func NewJWTManager(secret string) *JWTManagerImpl {
	return &JWTManagerImpl{secret: secret}
}

// Generate generates a JWT token for a user.
func (j *JWTManagerImpl) Generate(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// Verify validates a JWT token and returns the claims.
func (j *JWTManagerImpl) Verify(tokenStr string) (*models.AuthClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	userID, _ := claims["user_id"].(string)
	email, _ := claims["email"].(string)
	return &models.AuthClaims{UserID: userID, Email: email}, nil
}
