// Package models contains domain models for authentication and user management.
package models

// User represents a user in the system.
type User struct {
	ID       string
	Email    string
	Password string // hashed
}
