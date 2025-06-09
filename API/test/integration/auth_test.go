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

func TestAuthIntegration(t *testing.T) {
	ctx := context.Background()

	// Setup test container
	tc := SetupMongoDB(ctx, t)
	defer tc.TeardownMongoDB(ctx, t)

	// Setup test app
	app := SetupTestApp(tc)

	t.Run("should register and login users with different types", func(t *testing.T) {
		// Clean database before test
		tc.CleanDatabase(ctx, t)

		// Test data for different user types
		testUsers := []struct {
			name     string
			userType models.UserType
			email    string
			password string
			profile  models.UserProfile
		}{
			{
				name:     "Patient",
				userType: models.UserTypePatient,
				email:    "patient@test.com",
				password: "password123",
				profile: models.UserProfile{
					FirstName:   "João",
					LastName:    "Silva",
					Phone:       "+55-11-99999-9999",
					DateOfBirth: "1990-01-01",
					CPF:         "123.456.789-00",
				},
			},
			{
				name:     "Doctor",
				userType: models.UserTypeDoctor,
				email:    "doctor@test.com",
				password: "password123",
				profile: models.UserProfile{
					FirstName:  "Dr. Maria",
					LastName:   "Santos",
					Phone:      "+55-11-88888-8888",
					CRM:        "CRM-SP 123456",
					Speciality: "Cardiologia",
					Department: "Cardiologia",
				},
			},
			{
				name:     "Admin",
				userType: models.UserTypeAdmin,
				email:    "admin@test.com",
				password: "password123",
				profile: models.UserProfile{
					FirstName:  "Admin",
					LastName:   "Sistema",
					Phone:      "+55-11-77777-7777",
					Department: "TI",
				},
			},
		}

		for _, user := range testUsers {
			t.Run("Register_"+user.name, func(t *testing.T) {
				// Register user
				registerReq := models.RegisterRequest{
					Email:    user.email,
					Password: user.password,
					Type:     user.userType,
					Profile:  user.profile,
				}

				reqBody, err := json.Marshal(registerReq)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				app.Echo.ServeHTTP(rec, req)

				assert.Equal(t, http.StatusCreated, rec.Code)

				var registerResp models.RegisterResponse
				err = json.Unmarshal(rec.Body.Bytes(), &registerResp)
				require.NoError(t, err)

				assert.NotEmpty(t, registerResp.ID)
				assert.Equal(t, user.email, registerResp.Email)
				assert.Equal(t, user.userType, registerResp.Type)
				assert.Equal(t, user.profile.FirstName, registerResp.Profile.FirstName)
			})

			t.Run("Login_"+user.name, func(t *testing.T) {
				// Login user
				loginReq := models.LoginRequest{
					Email:    user.email,
					Password: user.password,
				}

				reqBody, err := json.Marshal(loginReq)
				require.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				app.Echo.ServeHTTP(rec, req)

				assert.Equal(t, http.StatusOK, rec.Code)

				var loginResp models.LoginResponse
				err = json.Unmarshal(rec.Body.Bytes(), &loginResp)
				require.NoError(t, err)

				assert.NotEmpty(t, loginResp.Token)

				// Validate token
				claims, err := app.JWTManager.Validate(loginResp.Token)
				require.NoError(t, err)
				assert.Equal(t, user.email, claims.Email)
				assert.Equal(t, user.userType, claims.UserType)
			})
		}
	})

	t.Run("should prevent duplicate user registration", func(t *testing.T) {
		// Clean database before test
		tc.CleanDatabase(ctx, t)

		registerReq := models.RegisterRequest{
			Email:    "duplicate@test.com",
			Password: "password123",
			Type:     models.UserTypePatient,
			Profile: models.UserProfile{
				FirstName: "Test",
				LastName:  "User",
				Phone:     "+55-11-99999-9999",
			},
		}

		reqBody, err := json.Marshal(registerReq)
		require.NoError(t, err)

		// First registration - should succeed
		req1 := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(reqBody))
		req1.Header.Set("Content-Type", "application/json")
		rec1 := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec1, req1)
		assert.Equal(t, http.StatusCreated, rec1.Code)

		// Second registration - should fail
		req2 := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(reqBody))
		req2.Header.Set("Content-Type", "application/json")
		rec2 := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec2, req2)
		assert.Equal(t, http.StatusConflict, rec2.Code)
	})

	t.Run("should reject invalid login credentials", func(t *testing.T) {
		// Clean database before test
		tc.CleanDatabase(ctx, t)

		// Register a user first
		registerReq := models.RegisterRequest{
			Email:    "valid@test.com",
			Password: "validpassword",
			Type:     models.UserTypePatient,
			Profile: models.UserProfile{
				FirstName: "Valid",
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
		assert.Equal(t, http.StatusCreated, rec.Code)

		// Try to login with wrong password
		loginReq := models.LoginRequest{
			Email:    "valid@test.com",
			Password: "wrongpassword",
		}

		reqBody, err = json.Marshal(loginReq)
		require.NoError(t, err)

		req = httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
