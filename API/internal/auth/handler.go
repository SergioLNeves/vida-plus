// Package auth provides authentication and authorization services.
package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	AuthService AuthService
}

func NewAuthHandler(authService AuthService) *AuthHandler {
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if req.Email == "" || req.Password == "" {
		return c.String(http.StatusBadRequest, "email and password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	user, err := h.AuthService.Register(c.Request().Context(), req.Email, string(hash))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.String(http.StatusConflict, "user already exists")
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, map[string]string{"id": user.ID, "email": user.Email})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if req.Email == "" || req.Password == "" {
		return c.String(http.StatusBadRequest, "email and password required")
	}
	token, err := h.AuthService.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.String(http.StatusUnauthorized, "invalid credentials")
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
