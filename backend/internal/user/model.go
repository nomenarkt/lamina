package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// DepartmentAccess defines a role and optional rank in a department.
type DepartmentAccess struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Rank string `json:"rank,omitempty"`
}

// Child represents dependent info for HR profiles.
type Child struct {
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
}

// InviteRequest is used when an admin invites a new user.
type InviteRequest struct {
	Email        string `json:"email"`
	UserType     string `json:"user_type"`               // "internal" or "external"
	DurationDays int    `json:"duration_days,omitempty"` // only for external users
}

// User represents a Lamina user.
type User struct {
	ID                    int64              `db:"id"`
	Email                 string             `db:"email"`
	PasswordHash          string             `db:"password_hash"`
	Role                  string             `db:"role"`
	Status                string             `db:"status"`
	ConfirmationToken     *string            `db:"confirmation_token"`
	CreatedAt             time.Time          `db:"created_at"`
	FullName              *string            `db:"full_name"`
	UserType              string             `db:"user_type"`  // admin, internal, external
	EmployeeID            *int               `db:"company_id"` // mapped from "employee_id"
	Phone                 *string            `db:"phone"`
	Address               *string            `db:"address"`
	ProfilePictureURL     *string            `db:"profile_picture_url" json:"profile_picture_url,omitempty"`
	Sex                   *string            `db:"sex" json:"sex,omitempty"` // male, female, non_binary
	Birthday              *string            `db:"birthday" json:"birthday,omitempty"`
	MaritalStatus         *bool              `db:"marital_status" json:"marital_status,omitempty"`
	SpouseName            *string            `db:"spouse_name" json:"spouse_name,omitempty"`
	HasChildren           *bool              `db:"has_children" json:"has_children,omitempty"`
	NumberOfChildren      *int               `db:"number_of_children" json:"number_of_children,omitempty"`
	Children              []Child            `json:"children,omitempty"` // not stored directly in DB unless serialized or in separate table
	NationalID            *string            `db:"national_id" json:"national_id,omitempty"`
	EmergencyContactName  *string            `db:"emergency_contact_name" json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string            `db:"emergency_contact_phone" json:"emergency_contact_phone,omitempty"`
	Departments           []DepartmentAccess `json:"departments"`
	Modules               []string           `json:"modules"`
	AccessExpiresAt       *time.Time         `db:"access_expires_at" json:"access_expires_at,omitempty"`
}

// Claims = JWT payload structure.
type Claims struct {
	UserID      int64              `json:"userID"`
	Email       string             `json:"email"`
	Departments []DepartmentAccess `json:"departments"`
	Modules     []string           `json:"modules"`
	Role        string             `json:"role"`
	jwt.RegisteredClaims
}

// UpdateProfileRequest represents the input required to update a user's profile.
type UpdateProfileRequest struct {
	FullName              string  `json:"full_name" binding:"required"`
	EmployeeID            *int    `json:"employee_id,omitempty"`
	Phone                 *string `json:"phone,omitempty"`
	Address               *string `json:"address,omitempty"`
	ProfilePictureURL     *string `json:"profile_picture_url,omitempty"`
	Sex                   *string `json:"sex,omitempty"`
	Birthday              *string `json:"birthday,omitempty"` // YYYY-MM-DD
	MaritalStatus         *bool   `json:"marital_status,omitempty"`
	SpouseName            *string `json:"spouse_name,omitempty"`
	HasChildren           *bool   `json:"has_children,omitempty"`
	NumberOfChildren      *int    `json:"number_of_children,omitempty"`
	Children              []Child `json:"children,omitempty"`
	NationalID            *string `json:"national_id,omitempty"`
	EmergencyContactName  *string `json:"emergency_contact_name,omitempty"`
	EmergencyContactPhone *string `json:"emergency_contact_phone,omitempty"`
}
