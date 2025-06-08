// Package handlers provides HTTP handlers for the API.
package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/vida-plus/api/models"
)

type AuthHandler struct {
	AuthService models.AuthService
}

func NewAuthHandler(authService models.AuthService) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
	}
}

func (h *AuthHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	logger := slog.With(
		slog.String("handler", "VoteHandler"),
		slog.String("func", "GetVoteByPollId"),
	)
	var req models.RegisterRequest

	if err := c.Bind(&req); err != nil {
		logger.Error("error during binding request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, models.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	if err := GetValidator().Struct(req); err != nil {
		logger.Error("error during validate struct values", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, models.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	user, err := h.AuthService.Register(ctx, req.Email, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			logger.Error("user already exists", slog.Any("error", err))
			return c.JSON(http.StatusConflict, models.NewAPIError(http.StatusConflict, err.Error()))
		}
		logger.Error("internal server error", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, models.NewAPIError(http.StatusInternalServerError, err.Error()))
	}

	return c.JSON(http.StatusCreated, models.RegisterResponse{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	logger := slog.With(
		slog.String("handler", "VoteHandler"),
		slog.String("func", "GetVoteByPollId"),
	)

	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		logger.Error("error during binding request", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, models.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	if err := GetValidator().Struct(req); err != nil {
		logger.Error("error during validate struct values", slog.Any("error", err))
		return c.JSON(http.StatusBadRequest, models.NewAPIError(http.StatusBadRequest, err.Error()))
	}

	token, err := h.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		logger.Error("error during login", slog.Any("error", err))
		return c.JSON(http.StatusUnauthorized, models.NewAPIError(http.StatusUnauthorized, err.Error()))
	}

	return c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
	})
}
