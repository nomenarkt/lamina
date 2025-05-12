// Package auth provides authentication and user credential management services.
package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/internal/user"
)

// Repository defines the methods implemented by the database layer.
type Repository interface {
	IsEmailExists(email string) (bool, error)
	CreateUser(ctx context.Context, companyID int, email string, hash string) (int64, error)
	FindByEmail(ctx context.Context, email string) (user.User, error)
}

// repositoryImpl handles database operations for user authentication.
type repositoryImpl struct {
	db *sqlx.DB
}

// NewAuthRepository returns a new Repository interface backed by PostgreSQL.
func NewAuthRepository(db *sqlx.DB) Repository {
	return &repositoryImpl{db: db}
}

// FindByEmail retrieves a user by email.
func (r *repositoryImpl) FindByEmail(ctx context.Context, email string) (user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, "SELECT id, email, password_hash, role FROM users WHERE email=$1", email)
	return u, err
}

// CreateUser inserts a new user into the database and returns the new user ID.
func (r *repositoryImpl) CreateUser(ctx context.Context, _ int, email string, hash string) (int64, error) {
	var id int64
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO users (email, password_hash, role, status, created_at)
		VALUES ($1, $2, 'user', 'pending', CURRENT_TIMESTAMP)
		RETURNING id
	`, email, hash).Scan(&id)

	return id, err
}

// IsEmailExists checks if a user with the given email already exists.
func (r *repositoryImpl) IsEmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
