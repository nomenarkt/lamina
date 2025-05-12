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
	if u.CompanyID != nil && *u.CompanyID > 0 {
		companyColumn = ", company_id"
		companyValue = ", :company_id"
	}

	query := `
		INSERT INTO users (email, password_hash, role, status, full_name, created_at` + companyColumn + `)
		VALUES (:email, :password_hash, :role, :status, :full_name, :created_at` + companyValue + `)
	`

	// Safe binding args
	args := map[string]interface{}{
		"email":         u.Email,
		"password_hash": u.PasswordHash,
		"role":          u.Role,
		"status":        u.Status,
		"full_name":     u.FullName,
		"created_at":    u.CreatedAt,
	}

	if u.CompanyID != nil && *u.CompanyID > 0 {
		args["company_id"] = u.CompanyID
	}

	_, err := r.db.NamedExecContext(ctx, query, args)
	return err
}
