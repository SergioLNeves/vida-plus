package handlers

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/models"
)

// AdminHandler handles admin-specific endpoints
type AdminHandler struct{}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetAllUsers godoc
// @Summary Get all users (Admin only)
// @Description Get list of all users in the system
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of all users"
// @Failure 401 {object} models.APIError "Unauthorized"
// @Failure 403 {object} models.APIError "Forbidden"
// @Router /admin/users [get]
func (h *AdminHandler) GetAllUsers(c echo.Context) error {
	logger := slog.With(
		slog.String("handler", "AdminHandler"),
		slog.String("func", "GetAllUsers"),
	)

	claims, err := models.GetAuthClaims(c.Get("claims"))
	if err != nil {
		logger.Error("missing or invalid claims", slog.Any("error", err))
		return c.JSON(http.StatusUnauthorized, models.NewAPIError(http.StatusUnauthorized, err.Error()))
	}

	// Simulação de lista de todos os usuários
	users := map[string]interface{}{
		"admin_id": claims.UserID,
		"users": []map[string]interface{}{
			{
				"id":         "user1",
				"email":      "joao@example.com",
				"type":       "patient",
				"status":     "active",
				"created_at": "2024-01-01",
			},
			{
				"id":         "user2",
				"email":      "dr.silva@example.com",
				"type":       "doctor",
				"status":     "active",
				"created_at": "2024-01-02",
			},
			{
				"id":         "user3",
				"email":      "nurse@example.com",
				"type":       "nurse",
				"status":     "active",
				"created_at": "2024-01-03",
			},
		},
	}

	logger.Info("admin accessed all users",
		slog.String("adminID", claims.UserID),
	)

	return c.JSON(http.StatusOK, users)
}

// GetSystemStats godoc
// @Summary Get system statistics (Admin only)
// @Description Get system usage statistics and metrics
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "System statistics"
// @Failure 401 {object} models.APIError "Unauthorized"
// @Failure 403 {object} models.APIError "Forbidden"
// @Router /admin/stats [get]
func (h *AdminHandler) GetSystemStats(c echo.Context) error {
	claims, err := models.GetAuthClaims(c.Get("claims"))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.NewAPIError(http.StatusUnauthorized, err.Error()))
	}

	// Simulação de estatísticas do sistema
	stats := map[string]interface{}{
		"total_users":        150,
		"total_patients":     120,
		"total_doctors":      15,
		"total_nurses":       10,
		"total_appointments": 500,
		"active_sessions":    25,
		"system_status":      "healthy",
	}

	slog.Info("admin accessed system stats",
		slog.String("adminID", claims.UserID),
	)

	return c.JSON(http.StatusOK, stats)
}
