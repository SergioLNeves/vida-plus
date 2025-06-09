// Package jwt provides JWT utilities for authentication.
package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vida-plus/api/internal/domain"
)

// JWTManagerImpl implements JWTManager interface.
type JWTManagerImpl struct {
	secret string
}

func NewJWTManager() domain.JWTManager {
	return &JWTManagerImpl{secret: "local-development-secret-key"} // Chave fixa para desenvolvimento local
}

// Generate generates a JWT token for a user.
func (j *JWTManagerImpl) Generate(user *domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"email":     user.Email,
		"user_type": user.Type,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// Validate validates a JWT token and returns the claims.
func (j *JWTManagerImpl) Validate(tokenStr string) (*domain.AuthClaims, error) {
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
	userTypeStr, _ := claims["user_type"].(string)
	userType := domain.UserType(userTypeStr)

	return &domain.AuthClaims{
		UserID:   userID,
		Email:    email,
		UserType: userType,
	}, nil
}

func (j *JWTManagerImpl) GenerateRefreshToken(user *domain.User) (string, error) {
	claims := &domain.AuthClaims{
		UserID: user.GetID(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 dias
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTManagerImpl) ValidateRefreshToken(token string) (*domain.AuthClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &domain.AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsedToken.Claims.(*domain.AuthClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
