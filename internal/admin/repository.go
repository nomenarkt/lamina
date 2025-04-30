package admin

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) CreateUser(ctx context.Context, email, passwordHash, role string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3)`,
		email, passwordHash, role,
	)
	return err
}
