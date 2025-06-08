// Package middleware provides HTTP middleware for authentication.
package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/internal/auth"
)

// AuthMiddleware checks for a valid JWT token.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: extract and validate JWT from Authorization header
		next.ServeHTTP(w, r)
	})
}

// JWTMiddleware validates JWT from Authorization header and sets user info in context.
func JWTMiddleware(jwtManager auth.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid token"})
			}
			token := strings.TrimPrefix(header, "Bearer ")
			claims, err := jwtManager.Verify(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			return next(c)
		}
	}
}
