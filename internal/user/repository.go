// Package user provides repository logic for accessing user-related data in the database.
package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Repo defines the contract for user-related database operations.
type Repo interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
	UpdateFullName(ctx context.Context, userID int64, fullName string) error
	UpdateUserProfile(ctx context.Context, userID int64, fullName string, companyID *int) error
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

// UpdateFullName changes the full name of the specified user.
func (r *Repository) UpdateFullName(ctx context.Context, userID int64, fullName string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET full_name = $1 WHERE id = $2`, fullName, userID)
	return err
}

// UpdateUserProfile updates both the full name and optionally the company ID of a user.
func (r *Repository) UpdateUserProfile(ctx context.Context, userID int64, fullName string, companyID *int) error {
	query := `UPDATE users SET full_name = :full_name`
	args := map[string]interface{}{
		"user_id":   userID,
		"full_name": fullName,
	}

	if companyID != nil {
		query += `, company_id = :company_id`
		args["company_id"] = *companyID
	}

	query += ` WHERE id = :user_id`

	_, err := r.db.NamedExecContext(ctx, query, args)
	return err
}
