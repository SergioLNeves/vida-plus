package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/internal/domain"
)

// RequirePermission creates a middleware that checks if user has specific permission
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := domain.GetAuthClaims(c.Get("claims"))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, domain.NewAPIError(
					http.StatusUnauthorized,
					"authentication required",
				))
			}

			// Buscar o usuário completo para verificar permissões
			// Por simplicidade, vamos usar apenas o tipo do usuário do token
			user := &domain.User{Type: claims.UserType}

			if !user.HasPermission(permission) {
				return c.JSON(http.StatusForbidden, domain.NewAPIError(
					http.StatusForbidden,
					"insufficient permissions",
				))
			}

			return next(c)
		}
	}
}

// RequireUserType creates a middleware that checks if user is of specific type
func RequireUserType(userTypes ...domain.UserType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := domain.GetAuthClaims(c.Get("claims"))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, domain.NewAPIError(
					http.StatusUnauthorized,
					"authentication required",
				))
			}

			// Verificar se o tipo do usuário está na lista permitida
			for _, allowedType := range userTypes {
				if claims.UserType == allowedType {
					return next(c)
				}
			}

			return c.JSON(http.StatusForbidden, domain.NewAPIError(
				http.StatusForbidden,
				"access denied for user type",
			))
		}
	}
}

// RequireAdmin creates a middleware that only allows admin users
func RequireAdmin() echo.MiddlewareFunc {
	return RequireUserType(domain.UserTypeAdmin)
}

// RequireMedicalStaff creates a middleware that allows doctors and nurses
func RequireMedicalStaff() echo.MiddlewareFunc {
	return RequireUserType(domain.UserTypeDoctor, domain.UserTypeNurse)
}

// RequireRole creates a middleware that checks if user has the required role
func RequireRole(requiredRole domain.UserType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := domain.GetAuthClaims(c.Get("claims"))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, domain.NewAPIError(
					http.StatusUnauthorized,
					"authentication required",
				))
			}

			if claims.UserType != requiredRole {
				return c.JSON(http.StatusForbidden, domain.NewAPIError(
					http.StatusForbidden,
					"insufficient role permissions",
				))
			}

			return next(c)
		}
	}
}
