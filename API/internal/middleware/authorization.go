package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/models"
)

// RequirePermission creates a middleware that checks if user has specific permission
func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := models.GetAuthClaims(c.Get("claims"))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.NewAPIError(
					http.StatusUnauthorized,
					"authentication required",
				))
			}

			// Buscar o usuário completo para verificar permissões
			// Por simplicidade, vamos usar apenas o tipo do usuário do token
			user := &models.User{Type: claims.UserType}

			if !user.HasPermission(permission) {
				return c.JSON(http.StatusForbidden, models.NewAPIError(
					http.StatusForbidden,
					"insufficient permissions",
				))
			}

			return next(c)
		}
	}
}

// RequireUserType creates a middleware that checks if user is of specific type
func RequireUserType(userTypes ...models.UserType) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := models.GetAuthClaims(c.Get("claims"))
			if err != nil {
				return c.JSON(http.StatusUnauthorized, models.NewAPIError(
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

			return c.JSON(http.StatusForbidden, models.NewAPIError(
				http.StatusForbidden,
				"access denied for user type",
			))
		}
	}
}

// RequireAdmin creates a middleware that only allows admin users
func RequireAdmin() echo.MiddlewareFunc {
	return RequireUserType(models.UserTypeAdmin)
}

// RequireMedicalStaff creates a middleware that allows doctors and nurses
func RequireMedicalStaff() echo.MiddlewareFunc {
	return RequireUserType(models.UserTypeDoctor, models.UserTypeNurse)
}
