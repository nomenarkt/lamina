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
