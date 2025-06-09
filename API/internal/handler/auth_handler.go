// Package handler provides HTTP handlers for the API.
package handler

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/vida-plus/api/internal/domain"
)

type AuthHandler struct {
	AuthService domain.AuthService
}

func NewAuthHandler(authService domain.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, password, type and profile
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "User registration data"
// @Success 201 {object} domain.RegisterResponse "User registered successfully"
// @Failure 400 {object} domain.APIError "Bad request"
// @Failure 409 {object} domain.APIError "User already exists"
// @Failure 500 {object} domain.APIError "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	logger := slog.With(
		slog.String("handler", "AuthHandler"),
		slog.String("func", "Register"),
	)
	var req domain.RegisterRequest

	if err := c.Bind(&req); err != nil {
		logger.Error("error during binding request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, domain.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	if err := GetValidator().Struct(req); err != nil {
		logger.Error("error during validate struct values", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, domain.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	user, err := h.AuthService.RegisterWithProfile(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			logger.Error("user already exists", slog.Any("error", err))
			return c.JSON(http.StatusConflict, domain.NewAPIError(http.StatusConflict, err.Error()))
		}
		logger.Error("internal server error", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, domain.NewAPIError(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusCreated, domain.RegisterResponse{
		ID:      user.ID,
		Email:   user.Email,
		Type:    user.Type,
		Profile: user.Profile,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "User login credentials"
// @Success 200 {object} domain.LoginResponse "Login successful"
// @Failure 400 {object} domain.APIError "Bad request"
// @Failure 401 {object} domain.APIError "Invalid credentials"
// @Failure 500 {object} domain.APIError "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	logger := slog.With(
		slog.String("handler", "AuthHandler"),
		slog.String("func", "Login"),
	)

	var req domain.LoginRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("error during binding request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, domain.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	if err := GetValidator().Struct(req); err != nil {
		logger.Error("error during validate struct values", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, domain.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	token, err := h.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.Error("error during login", slog.Any("error", err))
		return c.JSON(http.StatusUnauthorized, domain.NewAPIError(http.StatusUnauthorized, err.Error()))
	}

	return c.JSON(http.StatusOK, domain.LoginResponse{
		Token: token,
	})
}
