package integration

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckIntegration(t *testing.T) {
	ctx := context.Background()

	// Setup test container
	tc := SetupMongoDB(ctx, t)
	defer tc.TeardownMongoDB(ctx, t)

	// Setup test app
	app := SetupTestApp(tc)

	t.Run("should return healthy status when database is connected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "healthy")
	})

	t.Run("should be accessible without authentication", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		// Note: No Authorization header
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
