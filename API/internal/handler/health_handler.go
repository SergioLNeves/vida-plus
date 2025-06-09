package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	mongoClient *mongo.Client
}

func NewHealthHandler(mongoClient *mongo.Client) *HealthHandler {
	return &HealthHandler{
		mongoClient: mongoClient,
	}
}

// Check godoc
// @Summary Health check
// @Description Check the health status of the API and database connection
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Service is healthy"
// @Failure 503 {object} domain.APIError "Service unavailable"
// @Router /health [get]
func (h *HealthHandler) Check(c echo.Context) error {
	if err := h.mongoClient.Ping(c.Request().Context(), nil); err != nil {
		return c.JSON(http.StatusServiceUnavailable, domain.NewAPIError(
			http.StatusServiceUnavailable,
			"database health check failed",
		))
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
