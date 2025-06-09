package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/internal/domain"
)

// AdminHandler handles admin-specific endpoints
type AdminHandler struct {
	userRepo domain.UserRepository
}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler(userRepo domain.UserRepository) *AdminHandler {
	return &AdminHandler{
		userRepo: userRepo,
	}
}

// GetAllUsers godoc
// @Summary Get all users (Admin only)
// @Description Get list of all users in the system
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of all users"
// @Failure 401 {object} domain.APIError "Unauthorized"
// @Failure 403 {object} domain.APIError "Forbidden"
// @Router /admin/users [get]
func (h *AdminHandler) GetAllUsers(c echo.Context) error {
	logger := slog.With(
		slog.String("handler", "AdminHandler"),
		slog.String("func", "GetAllUsers"),
	)

	claims, err := domain.GetAuthClaims(c.Get("claims"))
	if err != nil {
		logger.Error("missing or invalid claims", slog.Any("error", err))
		return c.JSON(http.StatusUnauthorized, domain.NewAPIError(http.StatusUnauthorized, err.Error()))
	}

	// Busca real de usuários do banco de dados
	users, err := h.userRepo.GetAllUsers(c.Request().Context())
	if err != nil {
		logger.Error("failed to get all users", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, domain.NewAPIError(http.StatusInternalServerError, "failed to retrieve users"))
	}

	response := map[string]interface{}{
		"admin_id":    claims.UserID,
		"users":       users,
		"total_count": len(users),
	}

	logger.Info("admin accessed all users",
		slog.String("adminID", claims.UserID),
		slog.Int("userCount", len(users)),
	)

	return c.JSON(http.StatusOK, response)
}

// GetSystemStats godoc
// @Summary Get system statistics (Admin only)
// @Description Get system usage statistics and metrics
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "System statistics"
// @Failure 401 {object} domain.APIError "Unauthorized"
// @Failure 403 {object} domain.APIError "Forbidden"
// @Router /admin/stats [get]
func (h *AdminHandler) GetSystemStats(c echo.Context) error {
	logger := slog.With(
		slog.String("handler", "AdminHandler"),
		slog.String("func", "GetSystemStats"),
	)

	claims, err := domain.GetAuthClaims(c.Get("claims"))
	if err != nil {
		logger.Error("missing or invalid claims", slog.Any("error", err))
		return c.JSON(http.StatusUnauthorized, domain.NewAPIError(http.StatusUnauthorized, err.Error()))
	}

	// Get all users to calculate statistics
	users, err := h.userRepo.GetAllUsers(c.Request().Context())
	if err != nil {
		logger.Error("failed to get all users for stats", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, domain.NewAPIError(http.StatusInternalServerError, "failed to retrieve user statistics"))
	}

	// Count users by type
	userCounts := map[domain.UserType]int{
		domain.UserTypePatient:      0,
		domain.UserTypeDoctor:       0,
		domain.UserTypeNurse:        0,
		domain.UserTypeAdmin:        0,
		domain.UserTypeReceptionist: 0,
	}

	for _, user := range users {
		userCounts[user.Type]++
	}

	stats := map[string]interface{}{
		"admin_id":            claims.UserID,
		"total_users":         len(users),
		"total_patients":      userCounts[domain.UserTypePatient],
		"total_doctors":       userCounts[domain.UserTypeDoctor],
		"total_nurses":        userCounts[domain.UserTypeNurse],
		"total_admins":        userCounts[domain.UserTypeAdmin],
		"total_receptionists": userCounts[domain.UserTypeReceptionist],
	}

	logger.Info("admin accessed system stats",
		slog.String("adminID", claims.UserID),
		slog.Int("totalUsers", len(users)),
	)

	return c.JSON(http.StatusOK, stats)
}
