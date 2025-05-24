// Package admin provides the business logic for administrator operations.
package admin

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/auth"
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

// InviteUser invites a user by email and role, storing a token for confirmation/signup.
func (s *Service) InviteUser(ctx context.Context, req CreateUserRequest, _ string) error {
	if !utils.IsValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already registered")
	}

	// Determine user type based on email domain
	userType := "external"
	if strings.HasSuffix(strings.ToLower(req.Email), "@madagascarairlines.com") {
		userType = "internal"
	}

	// Default role to "user" if not specified
	role := req.Role
	if role == "" {
		role = "user"
	}

	newUser := &user.User{
		Email:        req.Email,
		PasswordHash: "",
		Status:       "pending",
		Role:         role,
		UserType:     userType,
		CreatedAt:    time.Now(),
	}

	// ✅ Set access_expires_at if external user with duration
	if userType == "external" && req.Duration != "" {
		dur, err := auth.ParseFlexibleDuration(req.Duration)
		if err != nil {
			return errors.New("invalid duration format (e.g., 2w, 1m, 90d)")
		}
		t := newUser.CreatedAt.Add(dur)
		newUser.AccessExpiresAt = &t
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		return err
	}

	userID, err := s.repo.FindUserIDByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return err
	}

	if err := s.repo.SetConfirmationToken(ctx, userID, token); err != nil {
		return err
	}

	return auth.SendConfirmationEmail(req.Email, token)
}
