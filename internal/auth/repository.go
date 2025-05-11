package auth

import (
	"context"

	"github.com/nomenarkt/lamina/internal/user"

	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (user.User, error) {
	var u user.User
	err := r.db.GetContext(ctx, &u, "SELECT id, email, password_hash, role FROM users WHERE email=$1", email)
	return u, err
}

func (r *AuthRepository) CreateUser(ctx context.Context, companyID int, email string, hash string) (int64, error) {
	var id int64
	err := r.db.QueryRowxContext(ctx, `
		INSERT INTO users (email, password_hash, role, status, created_at)
		VALUES ($1, $2, 'user', 'pending', CURRENT_TIMESTAMP)
		RETURNING id
	`, email, hash).Scan(&id)

	return id, err
}

func (r *AuthRepository) IsEmailExists(email string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
