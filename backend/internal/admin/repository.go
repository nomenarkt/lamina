// Package admin provides interfaces and implementations for managing admin operations.
package admin

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/internal/user"
)

// Repo defines the behavior for admin-related persistence logic.
type Repo interface {
	CreateUser(ctx context.Context, u *user.User) error
	IsEmailExists(email string) (bool, error)
	FindUserIDByEmail(ctx context.Context, email string) (int64, error)
	SetConfirmationToken(ctx context.Context, userID int64, token string) error
}

// Repository implements Repo using a SQL database.
type Repository struct {
	db *sqlx.DB
}

// NewAdminRepository creates a new Repository with the given database connection.
func NewAdminRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// CreateUser inserts a new user into the users table, optionally including company_id if present.
func (r *Repository) CreateUser(ctx context.Context, u *user.User) error {
	// Conditional additions
	companyColumn := ""
	companyValue := ""
	if u.EmployeeID != nil && *u.EmployeeID > 0 {
		companyColumn = ", company_id"
		companyValue = ", :company_id"
	}

	query := `
		INSERT INTO users (email, password_hash, role, status, full_name, created_at, user_type` + companyColumn + `)
		VALUES (:email, :password_hash, :role, :status, :full_name, :created_at, :user_type` + companyValue + `)
	`

	// Safe binding args
	args := map[string]interface{}{
		"email":         u.Email,
		"password_hash": u.PasswordHash,
		"role":          u.Role,
		"status":        u.Status,
		"full_name":     u.FullName,
		"created_at":    u.CreatedAt,
		"user_type":     u.UserType,
	}

	if u.EmployeeID != nil && *u.EmployeeID > 0 {
		args["company_id"] = u.EmployeeID
	}

	_, err := r.db.NamedExecContext(ctx, query, args)
	return err
}

// IsEmailExists checks if a user already exists with the given email address.
func (r *Repository) IsEmailExists(email string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM users WHERE email = $1", email)
	return count > 0, err
}

// FindUserIDByEmail retrieves the user ID associated with a specific email.
func (r *Repository) FindUserIDByEmail(ctx context.Context, email string) (int64, error) {
	var id int64
	err := r.db.GetContext(ctx, &id, "SELECT id FROM users WHERE email = $1", email)
	return id, err
}

// SetConfirmationToken assigns a confirmation token to the specified user ID.
func (r *Repository) SetConfirmationToken(ctx context.Context, userID int64, token string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET confirmation_token = $1
		WHERE id = $2
	`, token, userID)
	return err
}
