package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepo interface {
	FindByID(ctx context.Context, id int64) (*User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
	UpdateFullName(ctx context.Context, userID int64, fullName string) error
	UpdateUserProfile(ctx context.Context, userID int64, fullName string, companyID *int) error
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id=$1", id); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.SelectContext(ctx, &users, "SELECT id, email FROM users")
	return users, err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE email=$1", email); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) IsAdmin(ctx context.Context, id int64) (bool, error) {
	var role string
	err := r.db.GetContext(ctx, &role, "SELECT role FROM users WHERE id=$1", id)
	if err != nil {
		return false, err
	}
	return role == "admin", nil
}

func (r *UserRepository) UpdateFullName(ctx context.Context, userID int64, fullName string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE users SET full_name = $1 WHERE id = $2`, fullName, userID)
	return err
}

func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID int64, fullName string, companyID *int) error {
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
