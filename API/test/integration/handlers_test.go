package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vida-plus/api/models"
)

// TestHandlersIntegration tests all handler functionality with complete scenarios
func TestHandlersIntegration(t *testing.T) {
	ctx := context.Background()
	tc := SetupMongoDB(ctx, t)
	defer tc.Cleanup(ctx)

	app := SetupTestApp(tc)

	t.Run("Complete Patient Journey", func(t *testing.T) {
		// Step 1: Register a patient
		patientReq := models.RegisterRequest{
			Email:    "patient.journey@test.com",
			Password: "password123",
			Type:     models.UserTypePatient,
			Profile: models.UserProfile{
				FirstName:   "João",
				LastName:    "Silva",
				Phone:       "11999999999",
				DateOfBirth: "1990-01-01",
				CPF:         "123.456.789-00",
			},
		}

		body, _ := json.Marshal(patientReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var registerResp models.RegisterResponse
		err := json.Unmarshal(rec.Body.Bytes(), &registerResp)
		require.NoError(t, err)
		assert.Equal(t, models.UserTypePatient, registerResp.Type)

		// Step 2: Login as patient
		loginReq := models.LoginRequest{
			Email:    "patient.journey@test.com",
			Password: "password123",
		}

		body, _ = json.Marshal(loginReq)
		req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var loginResp models.LoginResponse
		err = json.Unmarshal(rec.Body.Bytes(), &loginResp)
		require.NoError(t, err)
		patientToken := loginResp.Token

		// Step 3: Access patient profile
		req = httptest.NewRequest(http.MethodGet, "/v1/patient/profile", nil)
		req.Header.Set("Authorization", "Bearer "+patientToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var profileResp struct {
			Profile models.UserProfile `json:"profile"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &profileResp)
		require.NoError(t, err)
		assert.Equal(t, "João", profileResp.Profile.FirstName)

		// Step 4: Access medical history
		req = httptest.NewRequest(http.MethodGet, "/v1/patient/medical-history", nil)
		req.Header.Set("Authorization", "Bearer "+patientToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Step 5: Try to access doctor-only endpoint (should fail)
		req = httptest.NewRequest(http.MethodGet, "/v1/doctor/patients", nil)
		req.Header.Set("Authorization", "Bearer "+patientToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("Complete Doctor Journey", func(t *testing.T) {
		// Step 1: Register a doctor
		doctorReq := models.RegisterRequest{
			Email:    "doctor.journey@test.com",
			Password: "password123",
			Type:     models.UserTypeDoctor,
			Profile: models.UserProfile{
				FirstName:  "Maria",
				LastName:   "Santos",
				Phone:      "11888888888",
				CRM:        "CRM12345",
				Speciality: "Cardiologia",
			},
		}

		body, _ := json.Marshal(doctorReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Step 2: Login as doctor
		loginReq := models.LoginRequest{
			Email:    "doctor.journey@test.com",
			Password: "password123",
		}

		body, _ = json.Marshal(loginReq)
		req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var loginResp struct {
			Token string `json:"token"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &loginResp)
		require.NoError(t, err)
		doctorToken := loginResp.Token

		// Step 3: Access patients list
		req = httptest.NewRequest(http.MethodGet, "/v1/doctor/patients", nil)
		req.Header.Set("Authorization", "Bearer "+doctorToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Step 4: Access appointments
		req = httptest.NewRequest(http.MethodGet, "/v1/doctor/appointments", nil)
		req.Header.Set("Authorization", "Bearer "+doctorToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Step 5: Try to access admin-only endpoint (should fail)
		req = httptest.NewRequest(http.MethodGet, "/v1/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+doctorToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("Complete Admin Journey", func(t *testing.T) {
		// Step 1: Register an admin
		adminReq := models.RegisterRequest{
			Email:    "admin.journey@test.com",
			Password: "password123",
			Type:     models.UserTypeAdmin,
			Profile: models.UserProfile{
				FirstName: "Admin",
				LastName:  "User",
				Phone:     "11777777777",
			},
		}

		body, _ := json.Marshal(adminReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Step 2: Login as admin
		loginReq := models.LoginRequest{
			Email:    "admin.journey@test.com",
			Password: "password123",
		}

		body, _ = json.Marshal(loginReq)
		req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var loginResp struct {
			Token string `json:"token"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &loginResp)
		require.NoError(t, err)
		adminToken := loginResp.Token

		// Step 3: Access user management
		req = httptest.NewRequest(http.MethodGet, "/v1/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var usersResp struct {
			Users []models.User `json:"users"`
			Total int           `json:"total"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &usersResp)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, usersResp.Total, 3) // At least the 3 users we created

		// Step 4: Access system stats
		req = httptest.NewRequest(http.MethodGet, "/v1/admin/stats", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var statsResp struct {
			TotalUsers    int `json:"total_users"`
			TotalPatients int `json:"total_patients"`
			TotalDoctors  int `json:"total_doctors"`
			TotalAdmins   int `json:"total_admins"`
		}
		err = json.Unmarshal(rec.Body.Bytes(), &statsResp)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, statsResp.TotalUsers, 3)
		assert.GreaterOrEqual(t, statsResp.TotalPatients, 1)
		assert.GreaterOrEqual(t, statsResp.TotalDoctors, 1)
		assert.GreaterOrEqual(t, statsResp.TotalAdmins, 1)

		// Step 5: Access shared doctor endpoint (admin has doctor permissions)
		req = httptest.NewRequest(http.MethodGet, "/v1/doctor/patients", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error Handling and Validation", func(t *testing.T) {
		// Test invalid registration data
		invalidReq := models.RegisterRequest{
			Email:    "invalid-email",
			Password: "123", // too short
			Type:     models.UserType("invalid_type"),
		}

		body, _ := json.Marshal(invalidReq)
		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Test invalid login
		invalidLogin := models.LoginRequest{
			Email:    "nonexistent@test.com",
			Password: "wrongpassword",
		}

		body, _ = json.Marshal(invalidLogin)
		req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Test accessing protected endpoint without token
		req = httptest.NewRequest(http.MethodGet, "/v1/patient/profile", nil)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// Test accessing protected endpoint with invalid token
		req = httptest.NewRequest(http.MethodGet, "/v1/patient/profile", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Health Check Comprehensive", func(t *testing.T) {
		// Test basic health endpoint
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var healthResp struct {
			Status string `json:"status"`
		}
		err := json.Unmarshal(rec.Body.Bytes(), &healthResp)
		require.NoError(t, err)
		assert.Equal(t, "healthy", healthResp.Status)

		// Test readiness endpoint
		req = httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Test liveness endpoint
		req = httptest.NewRequest(http.MethodGet, "/health/live", nil)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Cross-Role Interaction Scenarios", func(t *testing.T) {
		// Register multiple users for interaction testing
		users := []struct {
			email    string
			userType models.UserType
			name     string
		}{
			{"nurse1@test.com", models.UserTypeNurse, "Ana"},
			{"receptionist1@test.com", models.UserTypeReceptionist, "Carlos"},
		}

		tokens := make(map[string]string)

		for _, user := range users {
			// Register user
			regReq := models.RegisterRequest{
				Email:    user.email,
				Password: "password123",
				Type:     user.userType,
				Profile: models.UserProfile{
					FirstName: user.name,
					LastName:  "Test",
					Phone:     "11666666666",
				},
			}

			body, _ := json.Marshal(regReq)
			req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusCreated, rec.Code)

			// Login user
			loginReq := models.LoginRequest{
				Email:    user.email,
				Password: "password123",
			}

			body, _ = json.Marshal(loginReq)
			req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec = httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)

			var loginResp struct {
				Token string `json:"token"`
			}
			err := json.Unmarshal(rec.Body.Bytes(), &loginResp)
			require.NoError(t, err)
			tokens[string(user.userType)] = loginResp.Token
		}

		// Test nurse access patterns
		nurseToken := tokens[string(models.UserTypeNurse)]

		// Nurse should NOT have access to admin endpoints
		req := httptest.NewRequest(http.MethodGet, "/v1/admin/users", nil)
		req.Header.Set("Authorization", "Bearer "+nurseToken)
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)

		// Test receptionist access patterns
		receptionistToken := tokens[string(models.UserTypeReceptionist)]

		// Receptionist should NOT have access to doctor endpoints
		req = httptest.NewRequest(http.MethodGet, "/v1/doctor/patients", nil)
		req.Header.Set("Authorization", "Bearer "+receptionistToken)
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
