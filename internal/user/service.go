package user

import (
	"context"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetMe(ctx context.Context, id int64) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) GetProfile(ctx context.Context, id int64) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) FindAll(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	return s.repo.IsAdmin(ctx, userID)
}
