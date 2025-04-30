package user

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := r.db.GetContext(ctx, &u, "SELECT id, email FROM users WHERE id=$1", id)
	return u, err
}

func (r *UserRepository) FindAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.SelectContext(ctx, &users, "SELECT id, email FROM users")
	return users, err
}
