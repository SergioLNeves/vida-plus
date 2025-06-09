package healthcheck

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// HealthCheckService interface defines health check operations
type HealthCheckService interface {
	Status(ctx context.Context) *HealthCheckResultItem
}

// HealthCheckResultItem represents the result of a health check
type HealthCheckResultItem struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	WorkingStatus    = "working"
	NotWorkingStatus = "not_working"
)

// HealthCheckServiceImpl implements HealthCheckService
type HealthCheckServiceImpl struct {
	mongoClient *mongo.Client
}

// NewHealthCheckService creates a new HealthCheckService
func NewHealthCheckService(mongoClient *mongo.Client) HealthCheckService {
	return &HealthCheckServiceImpl{mongoClient: mongoClient}
}

// Status checks the health of the database connection
func (h *HealthCheckServiceImpl) Status(ctx context.Context) *HealthCheckResultItem {
	if h.mongoClient == nil {
		return &HealthCheckResultItem{
			Status: NotWorkingStatus,
			Error:  "MongoDB client not initialized",
		}
	}

	if err := h.mongoClient.Ping(ctx, nil); err != nil {
		return &HealthCheckResultItem{
			Status: NotWorkingStatus,
			Error:  "Failed to ping MongoDB: " + err.Error(),
		}
	}

	return &HealthCheckResultItem{
		Status: WorkingStatus,
		Error:  "",
	}
}
