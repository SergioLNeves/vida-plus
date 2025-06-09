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

// TestCoreIntegration tests the core functionality: authentication with user differentiation
func TestCoreIntegration(t *testing.T) {
	ctx := context.Background()
	tc := SetupMongoDB(ctx, t)
	defer tc.Cleanup(ctx)

	app := SetupTestApp(tc)

	t.Run("Complete Authentication and User Differentiation Flow", func(t *testing.T) {
		// Test different user types
		userTypes := []struct {
			name     string
			userType models.UserType
			email    string
		}{
			{"Patient", models.UserTypePatient, "patient@test.com"},
			{"Doctor", models.UserTypeDoctor, "doctor@test.com"},
			{"Admin", models.UserTypeAdmin, "admin@test.com"},
		}

		tokens := make(map[models.UserType]string)

		// Step 1: Register all user types
		for _, user := range userTypes {
			t.Run("Register_"+user.name, func(t *testing.T) {
				registerReq := models.RegisterRequest{
					Email:    user.email,
					Password: "password123",
					Type:     user.userType,
					Profile: models.UserProfile{
						FirstName: user.name,
						LastName:  "Test User",
						Phone:     "123456789",
					},
				}

				body, _ := json.Marshal(registerReq)
				req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				app.Echo.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusCreated, rec.Code)

				var registerResp models.RegisterResponse
				err := json.Unmarshal(rec.Body.Bytes(), &registerResp)
				require.NoError(t, err)
				assert.Equal(t, user.userType, registerResp.Type)
				assert.Equal(t, user.email, registerResp.Email)
			})
		}

		// Step 2: Login all users and get tokens
		for _, user := range userTypes {
			t.Run("Login_"+user.name, func(t *testing.T) {
				loginReq := models.LoginRequest{
					Email:    user.email,
					Password: "password123",
				}

				body, _ := json.Marshal(loginReq)
				req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				app.Echo.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusOK, rec.Code)

				var loginResp models.LoginResponse
				err := json.Unmarshal(rec.Body.Bytes(), &loginResp)
				require.NoError(t, err)
				assert.NotEmpty(t, loginResp.Token)

				tokens[user.userType] = loginResp.Token
			})
		}

		// Step 3: Test protected routes with different user types
		for _, user := range userTypes {
			t.Run("AccessProtectedRoute_"+user.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, "/v1/profile", nil)
				req.Header.Set("Authorization", "Bearer "+tokens[user.userType])
				rec := httptest.NewRecorder()

				app.Echo.ServeHTTP(rec, req)
				assert.Equal(t, http.StatusOK, rec.Code)

				var response map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.Equal(t, string(user.userType), response["type"])
				assert.Equal(t, user.email, response["email"])
				assert.Contains(t, response["message"], string(user.userType))
			})
		}
	})

	t.Run("Authorization Errors", func(t *testing.T) {
		t.Run("should_deny_access_without_token", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})

		t.Run("should_deny_access_with_invalid_token", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
			req.Header.Set("Authorization", "Bearer invalid_token")
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	})
}
