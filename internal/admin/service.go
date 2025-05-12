package admin

import (
	"context"
	"time"

	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/user"
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

	newUser := &user.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       "active",
		CreatedAt:    time.Now(),
	}

	return s.repo.CreateUser(ctx, newUser)
}
