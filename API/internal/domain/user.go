// Package models contains domain models for authentication and user management.
package domain

import (
	"context"
	"time"
)

// UserType represents the type of user in the system
type UserType string

const (
	UserTypePatient      UserType = "patient"      // Patient
	UserTypeDoctor       UserType = "doctor"       // Doctor
	UserTypeNurse        UserType = "nurse"        // Nurse
	UserTypeAdmin        UserType = "admin"        // Administrator
	UserTypeReceptionist UserType = "receptionist" // Receptionist
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusPending  UserStatus = "pending"
	UserStatusBlocked  UserStatus = "blocked"
)

// User represents a user in the system.
type User struct {
	ID        string      `bson:"_id" json:"id"`
	Email     string      `bson:"email" json:"email"`
	Password  string      `bson:"password" json:"-"` // hashed, not exposed in JSON
	Type      UserType    `bson:"type" json:"type"`
	Status    UserStatus  `bson:"status" json:"status"`
	Profile   UserProfile `bson:"profile" json:"profile"`
	CreatedAt time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time   `bson:"updated_at" json:"updated_at"`
}

// UserProfile contains profile information for all user types
type UserProfile struct {
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	Phone       string `bson:"phone" json:"phone"`
	DateOfBirth string `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"` // For patients
	CPF         string `bson:"cpf,omitempty" json:"cpf,omitempty"`
	CRM         string `bson:"crm,omitempty" json:"crm,omitempty"`               // For doctors
	COREN       string `bson:"coren,omitempty" json:"coren,omitempty"`           // For nurses
	Speciality  string `bson:"speciality,omitempty" json:"speciality,omitempty"` // For doctors
	Department  string `bson:"department,omitempty" json:"department,omitempty"` // For staff
}

// HasPermission checks if user has permission for a specific action
func (u *User) HasPermission(permission string) bool {
	switch u.Type {
	case UserTypeAdmin:
		return true // Admin tem acesso total
	case UserTypeDoctor:
		return permission == "view_patients" || permission == "manage_appointments" || permission == "view_medical_records"
	case UserTypeNurse:
		return permission == "view_patients" || permission == "view_basic_records"
	case UserTypeReceptionist:
		return permission == "manage_appointments" || permission == "view_patients"
	case UserTypePatient:
		return permission == "view_own_records" || permission == "manage_own_appointments"
	default:
		return false
	}
}

// IsActive checks if user account is active
func (u *User) IsActive() bool {
	return u.Status == UserStatusActive
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.Profile.FirstName + " " + u.Profile.LastName
}

// GetID returns the ID of the user
func (u *User) GetID() string {
	return u.ID
}

// UserStore defines user management methods.
type UserStore interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, user *User) error
}
