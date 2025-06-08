// Package main is the entry point for the API server.
package main

import (
	"github.com/labstack/echo/v4"
	"github.com/vida-plus/api/internal/auth"
	"github.com/vida-plus/api/internal/middleware"
	"github.com/vida-plus/api/internal/user"
	"github.com/vida-plus/api/pkg"
)

func ProtectedHandler(c echo.Context) error {
	userID := c.Get("userID")
	email := c.Get("email")
	return c.JSON(200, map[string]interface{}{"userID": userID, "email": email})
}

func RegisterHandler(authHandler *auth.AuthHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		return authHandler.Register(c)
	}
}

func LoginHandler(authHandler *auth.AuthHandler) echo.HandlerFunc {
	return func(c echo.Context) error {
		return authHandler.Login(c)
	}
}

func main() {
	userService := user.NewUserService()
	jwtManager := pkg.NewJWTManager("secret")
	e := echo.New()
	authService := auth.NewAuthService(userService, jwtManager)
	authHandler := auth.NewAuthHandler(authService)

	e.POST("/register", RegisterHandler(authHandler))
	e.POST("/login", LoginHandler(authHandler))
	e.GET("/protected", ProtectedHandler, middleware.JWTMiddleware(jwtManager))

	e.Logger.Fatal(e.Start(":8080"))
}
