// Package handler contains HTTP handlers for the API.
package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/internal/domain"
)

// ProtectedHandler struct holds handler dependencies
type ProtectedHandler struct{}

// NewProtectedHandler creates a new instance of ProtectedHandler
func NewProtectedHandler() *ProtectedHandler {
	return &ProtectedHandler{}
}

// GetProtectedInfo handles requests to protected endpoints.
// GetProtectedInfo godoc
// @Summary Get protected information
// @Description Get protected information that requires authentication
// @Tags protected
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Protected information"
// @Failure 401 {object} domain.APIError "Unauthorized"
// @Failure 500 {object} domain.APIError "Internal server error"
// @Router /protected [get]
func (h *ProtectedHandler) GetProtectedInfo(c echo.Context) error {
	logger := slog.With(
		slog.String("handler", "ProtectedHandler"),
		slog.String("func", "GetProtectedInfo"),
	)

	claims, err := domain.GetAuthClaims(c.Get("claims"))
	if err != nil {
		logger.Error("missing or invalid claims in context", slog.Any("error", err))
		return c.JSON(http.StatusUnauthorized, domain.NewAPIError(http.StatusUnauthorized, err.Error()))
	}

	logger.Info("protected info accessed",
		slog.String("userID", claims.UserID),
		slog.String("email", claims.Email),
	)

	return c.JSON(http.StatusOK, claims)
}

// RegisterRoutes registers protected routes in the Echo instance.
func RegisterProtectedRoutes(e *echo.Echo, jwtMiddleware echo.MiddlewareFunc) {
	handler := NewProtectedHandler()
	e.GET("/protected", handler.GetProtectedInfo, jwtMiddleware)
}
