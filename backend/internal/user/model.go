package user

import "time"

// DepartmentAccess defines a role and optional rank in a department.
type DepartmentAccess struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Rank string `json:"rank,omitempty"`
}

// User represents a Lamina user.
type User struct {
	ID                int64              `db:"id"`
	Email             string             `db:"email"`
	PasswordHash      string             `db:"password_hash"`
	Role              string             `db:"role"`
	Status            string             `db:"status"`
	ConfirmationToken *string            `db:"confirmation_token"`
	CreatedAt         time.Time          `db:"created_at"`
	FullName          *string            `db:"full_name"`
	UserType          string             `db:"user_type"`  // admin, internal, external
	EmployeeID        *int               `db:"company_id"` // mapped from "employee_id"
	Phone             *string            `db:"phone"`
	Address           *string            `db:"address"`
	Departments       []DepartmentAccess `json:"departments"`
	Modules           []string           `json:"modules"`
}

// Claims = JWT payload structure.
type Claims struct {
	UserID      int64              `json:"userID"`
	Email       string             `json:"email"`
	Departments []DepartmentAccess `json:"departments"`
	Modules     []string           `json:"modules"`
}

// UpdateProfileRequest represents the input required to update a user's profile.
type UpdateProfileRequest struct {
	FullName   string  `json:"full_name" binding:"required"`
	EmployeeID *int    `json:"employee_id,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Address    *string `json:"address,omitempty"`
}
