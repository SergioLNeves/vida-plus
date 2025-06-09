// Package models contains request and response models for the API.
package models

// RegisterRequest represents the request structure for user registration.
type RegisterRequest struct {
	Email    string      `json:"email" validate:"required,email" example:"user@example.com"`
	Password string      `json:"password" validate:"required,min=8,max=128" example:"mypassword123"`
	Type     UserType    `json:"type" validate:"required" example:"patient"`
	Profile  UserProfile `json:"profile" validate:"required"`
}

// LoginRequest represents the request structure for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=1" example:"mypassword123"`
}

// RegisterResponse represents the response structure for user registration.
type RegisterResponse struct {
	ID      string      `json:"id" example:"user123"`
	Email   string      `json:"email" example:"user@example.com"`
	Type    UserType    `json:"type" example:"patient"`
	Profile UserProfile `json:"profile"`
}

// LoginResponse represents the response structure for user login.
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ErrorResponse represents the error response structure.
type ErrorResponse struct {
	Error   string            `json:"error" example:"validation failed"`
	Details map[string]string `json:"details,omitempty"`
}
