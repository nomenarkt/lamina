package user

import (
	"context"
)

type UserService struct {
	repo *UserRepository
}

func NewUserService(r *UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) GetProfile(ctx context.Context, userID int64) (User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}
