package admin

import (
	"context"

	"github.com/nomenarkt/lamina/common/utils"
)

type AdminService struct {
	repo   AdminRepo
	hasher utils.PasswordHasher
}

func NewAdminService(repo AdminRepo, hasher utils.PasswordHasher) *AdminService {
	return &AdminService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *AdminService) CreateUser(ctx context.Context, req CreateUserRequest, createdBy string) error {
	hashedPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(ctx, req.Email, hashedPassword, req.Role)
}
