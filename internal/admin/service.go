package admin

import (
	"context"
)

type AdminService struct {
	repo *AdminRepository
}

func NewAdminService(r *AdminRepository) *AdminService {
	return &AdminService{repo: r}
}

func (s *AdminService) CreateUser(ctx context.Context, req CreateUserRequest, hashedPassword string) error {
	return s.repo.CreateUser(ctx, req.Email, hashedPassword, req.Role)
}
