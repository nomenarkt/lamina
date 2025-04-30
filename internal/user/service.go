package user

import (
	"context"
)

type UserRepo interface {
	FindByID(ctx context.Context, id int64) (User, error)
	FindAll(ctx context.Context) ([]User, error)
}

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) GetMe(ctx context.Context, userID int64) (User, error) {
	return s.repo.FindByID(ctx, userID)
}
