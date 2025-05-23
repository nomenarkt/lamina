// Package auth provides authentication and user credential management services.
package auth

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/internal/user"
)

// Repository defines the methods implemented by the database layer.
type Repository interface {
	FindByEmail(ctx context.Context, email string) (user.User, error)
	CreateUser(ctx context.Context, companyID int, email string, hash string) (int64, error)
	CreateUserWithType(ctx context.Context, companyID *int, email, hash, userType string) (int64, error)
	IsEmailExists(email string) (bool, error)
	FindByConfirmationToken(ctx context.Context, token string) (user.User, error)
	MarkUserConfirmed(ctx context.Context, id int64) error
	SetConfirmationToken(ctx context.Context, userID int64, token string) error
	UpdatePasswordAndActivate(ctx context.Context, userID int64, hashed string) error
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
	err := r.db.GetContext(ctx, &u, `
		SELECT id, email, password_hash, role, status
		FROM users WHERE email=$1
	`, email)
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

func (r *repositoryImpl) CreateUserWithType(ctx context.Context, companyID *int, email, hash, userType string) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users (company_id, email, password_hash, role, status, user_type)
		 VALUES ($1, $2, $3, 'user', 'pending', $4)
		 RETURNING id`, companyID, email, hash, userType).Scan(&id)
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

func (r *repositoryImpl) FindByConfirmationToken(ctx context.Context, token string) (user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, `
		SELECT id, email, status, created_at FROM users
		WHERE confirmation_token = $1
	`, token)
	return u, err
}

func (r *repositoryImpl) MarkUserConfirmed(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET status = 'active', confirmation_token = NULL
		WHERE id = $1
	`, id)
	return err
}

func (r *repositoryImpl) SetConfirmationToken(ctx context.Context, userID int64, token string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET confirmation_token = $1 WHERE id = $2
	`, token, userID)
	return err
}

func (r *repositoryImpl) UpdatePasswordAndActivate(ctx context.Context, userID int64, hashed string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users
		SET password_hash = $1,
		    status = 'active',
		    confirmation_token = NULL
		WHERE id = $2
	`, hashed, userID)
	return err
}
