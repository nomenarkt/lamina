// Package user provides repository logic for accessing user-related data in the database.
package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Repo defines the contract for user-related database operations.
type Repo interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
	UpdateUserProfile(ctx context.Context, userID int64, fullName string, employeeID *int, phone, address *string) error
	MarkUserActive(ctx context.Context, id int64) error
	DeleteExpiredPendingUsers(ctx context.Context) error
}

// Repository implements the Repo interface using sqlx for DB interaction.
type Repository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new instance of Repository.
func NewUserRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// FindByID retrieves a user by their ID.
func (r *Repository) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll returns a list of users with basic fields.
func (r *Repository) FindAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.SelectContext(ctx, &users, "SELECT id, email FROM users")
	return users, err
}

// FindByEmail looks up a user by their email address.
func (r *Repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		return nil, err
	}
	return &user, nil
}

// IsAdmin checks if a user has an "admin" role.
func (r *Repository) IsAdmin(ctx context.Context, id int64) (bool, error) {
	var role string
	err := r.db.GetContext(ctx, &role, "SELECT role FROM users WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}

// UpdateUserProfile updates the user's full name and optionally company ID, phone, and address.
func (r *Repository) UpdateUserProfile(ctx context.Context, userID int64, fullName string, employeeID *int, phone, address *string) error {
	var (
		setClauses []string
		args       = map[string]interface{}{
			"user_id": userID,
		}
	)

	// Required field
	setClauses = append(setClauses, "full_name = :full_name")
	args["full_name"] = fullName

	// Optional fields
	if employeeID != nil {
		setClauses = append(setClauses, "company_id = :company_id")
		args["company_id"] = *employeeID
	}
	if phone != nil {
		setClauses = append(setClauses, "phone = :phone")
		args["phone"] = *phone
	}
	if address != nil {
		setClauses = append(setClauses, "address = :address")
		args["address"] = *address
	}

	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = :user_id
	`, strings.Join(setClauses, ", "))

	_, err := r.db.NamedExecContext(ctx, query, args)
	return err
}

// MarkUserActive sets status = 'active' for a confirmed user.
func (r *Repository) MarkUserActive(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET status = 'active' WHERE id = $1`, id)
	return err
}

// DeleteExpiredPendingUsers removes pending users who didn't confirm within 24h.
func (r *Repository) DeleteExpiredPendingUsers(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM users
		WHERE status = 'pending'
		  AND created_at < NOW() - interval '24 hours'
	`)
	return err
}
