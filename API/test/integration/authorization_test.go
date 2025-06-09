package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vida-plus/api/models"
)

func TestAuthorizationIntegration(t *testing.T) {
	ctx := context.Background()

	// Setup test container
	tc := SetupMongoDB(ctx, t)
	defer tc.TeardownMongoDB(ctx, t)

	// Setup test app
	app := SetupTestApp(tc)

	// Helper function to register and login a user
	registerAndLogin := func(t *testing.T, userType models.UserType, email, password string) string {
		t.Helper()

		// Register user
		registerReq := models.RegisterRequest{
			Email:    email,
			Password: password,
			Type:     userType,
			Profile: models.UserProfile{
				FirstName: "Test",
				LastName:  "User",
				Phone:     "+55-11-99999-9999",
			},
		}

		reqBody, err := json.Marshal(registerReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)

		// Login user
		loginReq := models.LoginRequest{
			Email:    email,
			Password: password,
		}

		reqBody, err = json.Marshal(loginReq)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var loginResp models.LoginResponse
		err = json.Unmarshal(rec.Body.Bytes(), &loginResp)
		require.NoError(t, err)

		return loginResp.Token
	}

	t.Run("Patient Access Control", func(t *testing.T) {
		// Clean database before test
		tc.CleanDatabase(ctx, t)

		patientToken := registerAndLogin(t, models.UserTypePatient, "patient@test.com", "password123")

		t.Run("should allow patient to access basic profile", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/profile", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", patientToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("should allow patient to access protected endpoint", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", patientToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("should deny patient access to admin routes", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/admin/users", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", patientToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusForbidden, rec.Code)
		})
	})

	t.Run("Doctor Access Control", func(t *testing.T) {
		// Clean database before test
		tc.CleanDatabase(ctx, t)

		doctorToken := registerAndLogin(t, models.UserTypeDoctor, "doctor@test.com", "password123")

		t.Run("should allow doctor to access basic profile", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/profile", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", doctorToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("should allow doctor to access protected endpoint", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", doctorToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("should deny doctor access to admin routes", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/admin/users", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", doctorToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusForbidden, rec.Code)
		})
	})

	t.Run("Admin Access Control", func(t *testing.T) {
		// Clean database before test
		tc.CleanDatabase(ctx, t)

		adminToken := registerAndLogin(t, models.UserTypeAdmin, "admin@test.com", "password123")

		t.Run("should allow admin to access user management", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/admin/users", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("should allow admin to access system stats", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/admin/stats", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("should allow admin to access basic profile", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/profile", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})

		t.Run("should allow admin to access protected endpoint", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", adminToken))
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		})
	})

	t.Run("Unauthorized Access", func(t *testing.T) {
		t.Run("should deny access without token", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/profile", nil)
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})

		t.Run("should deny access with invalid token", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/profile", nil)
			req.Header.Set("Authorization", "Bearer invalid-token")
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	})
}
