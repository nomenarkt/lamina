// Package admin provides the business logic for administrator operations.
package admin

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/user"
)

// Service provides operations for administrative tasks.
type Service struct {
	repo   Repo
	hasher utils.PasswordHasher
}

// NewAdminService creates a new instance of Service.
func NewAdminService(repo Repo, hasher utils.PasswordHasher) *Service {
	return &Service{
		repo:   repo,
		hasher: hasher,
	}
}

// CreateUser creates a new user based on the admin request.
// The 'createdBy' parameter is currently unused, but reserved for audit logging or future features.
func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest, _ string) error {
	if !strings.HasSuffix(strings.ToLower(req.Email), "@madagascarairlines.com") {
		return errors.New("only corporate emails allowed for admin-created users")
	}

	hashedPassword, err := s.hasher.HashPassword(req.Password)
	if err != nil {
		return err
	}

	newUser := &user.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Status:       "pending",
		Role:         "", // default role is empty unless assigned later
		CreatedAt:    time.Now(),
	}

	return s.repo.CreateUser(ctx, newUser)
}
