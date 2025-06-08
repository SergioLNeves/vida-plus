// Package models contains domain models for authentication and authorization.
package models

// AuthClaims represents JWT claims for authentication.
type AuthClaims struct {
	UserID string
	Email  string
}
