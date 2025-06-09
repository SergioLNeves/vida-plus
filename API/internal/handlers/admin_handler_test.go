package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockmodels "github.com/vida-plus/api/mocks"
	"github.com/vida-plus/api/models"
)

func TestAdminHandler_GetAllUsers(t *testing.T) {
	t.Run("should return all users when request succeeds", func(t *testing.T) {
		// Setup
		mockUserRepo := &mockmodels.UserRepository{}
		handler := NewAdminHandler(mockUserRepo)

		// Test data
		users := []*models.User{
			{
				ID:     "user-1",
				Email:  "joao@example.com",
				Type:   models.UserTypePatient,
				Status: models.UserStatusActive,
				Profile: models.UserProfile{
					FirstName: "João",
					LastName:  "Silva",
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:     "user-2",
				Email:  "maria@example.com",
				Type:   models.UserTypeAdmin,
				Status: models.UserStatusActive,
				Profile: models.UserProfile{
					FirstName: "Maria",
					LastName:  "Santos",
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Mock expectations
		mockUserRepo.On("GetAllUsers", mock.Anything).Return(users, nil)

		// Setup request
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add user context (simulating authenticated admin)
		claims := &models.AuthClaims{
			UserID:   "admin-user-id",
			Email:    "admin@example.com",
			UserType: models.UserTypeAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "test-user",
			},
		}
		c.Set("claims", claims)

		// Execute
		err := handler.GetAllUsers(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "admin-user-id", response["admin_id"])
		assert.Contains(t, response, "users")
		assert.Contains(t, response, "total_count")
		assert.Equal(t, float64(2), response["total_count"])

		// Verify mock expectations
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockUserRepo := &mockmodels.UserRepository{}
		handler := NewAdminHandler(mockUserRepo)

		// Mock expectations
		mockUserRepo.On("GetAllUsers", mock.Anything).Return(nil, assert.AnError)

		// Setup request
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add user context
		claims := &models.AuthClaims{
			UserID:   "admin-user-id",
			Email:    "admin@example.com",
			UserType: models.UserTypeAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "test-user",
			},
		}
		c.Set("claims", claims)

		// Execute
		err := handler.GetAllUsers(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(500), response["status"])
		assert.Contains(t, response, "title")

		// Verify mock expectations
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAdminHandler_GetSystemStats(t *testing.T) {
	t.Run("should return system stats successfully", func(t *testing.T) {
		// Setup
		mockUserRepo := &mockmodels.UserRepository{}
		handler := NewAdminHandler(mockUserRepo)

		// Test data - create users of different types
		users := []*models.User{
			{
				ID:     "user-1",
				Email:  "patient@example.com",
				Type:   models.UserTypePatient,
				Status: models.UserStatusActive,
			},
			{
				ID:     "user-2",
				Email:  "doctor@example.com",
				Type:   models.UserTypeDoctor,
				Status: models.UserStatusActive,
			},
			{
				ID:     "user-3",
				Email:  "admin@example.com",
				Type:   models.UserTypeAdmin,
				Status: models.UserStatusActive,
			},
			{
				ID:     "user-4",
				Email:  "nurse@example.com",
				Type:   models.UserTypeNurse,
				Status: models.UserStatusActive,
			},
		}

		// Mock expectations
		mockUserRepo.On("GetAllUsers", mock.Anything).Return(users, nil)

		// Setup request
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add user context
		claims := &models.AuthClaims{
			UserID:   "admin-user-id",
			Email:    "admin@example.com",
			UserType: models.UserTypeAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "test-user",
			},
		}
		c.Set("claims", claims)

		// Execute
		err := handler.GetSystemStats(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "admin-user-id", response["admin_id"])
		assert.Equal(t, float64(4), response["total_users"])
		assert.Equal(t, float64(1), response["total_patients"])
		assert.Equal(t, float64(1), response["total_doctors"])
		assert.Equal(t, float64(1), response["total_admins"])
		assert.Equal(t, float64(1), response["total_nurses"])
		assert.Equal(t, float64(0), response["total_receptionists"])

		// Verify mock expectations
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		// Setup
		mockUserRepo := &mockmodels.UserRepository{}
		handler := NewAdminHandler(mockUserRepo)

		// Mock expectations
		mockUserRepo.On("GetAllUsers", mock.Anything).Return(nil, assert.AnError)

		// Setup request
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/stats", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Add user context
		claims := &models.AuthClaims{
			UserID:   "admin-user-id",
			Email:    "admin@example.com",
			UserType: models.UserTypeAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "test-user",
			},
		}
		c.Set("claims", claims)

		// Execute
		err := handler.GetSystemStats(c)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(500), response["status"])
		assert.Contains(t, response, "title")

		// Verify mock expectations
		mockUserRepo.AssertExpectations(t)
	})
}
