package admin

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nomenarkt/lamina/internal/user"
)

type AdminRepo interface {
	CreateUser(ctx context.Context, u *user.User) error
}

type AdminRepository struct {
	db *sqlx.DB
}

func NewAdminRepository(db *sqlx.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) CreateUser(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (company_id, email, password_hash, role, status, created_at)
		VALUES (:company_id, :email, :password_hash, :role, :status, :created_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, u)
	return err
}
