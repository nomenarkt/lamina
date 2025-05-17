// Package user implements business logic for user-related operations.
package user

import (
	"context"
	"errors"
	"strings"
)

// Service provides methods for working with user entities.
type Service struct {
	repo Repo
}

// NewUserService creates a new instance of Service.
func NewUserService(repo Repo) *Service {
	return &Service{repo: repo}
}

// GetMe retrieves the user associated with the given ID.
func (s *Service) GetMe(ctx context.Context, id int64) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

// FindAll retrieves a list of all users.
func (s *Service) FindAll(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

// ListUsers returns all users. This is similar to FindAll.
func (s *Service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

// UpdateUserProfile updates a user's full name and optionally their company ID,
// depending on whether the user's email belongs to the internal domain.
func (s *Service) UpdateUserProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// ⛔ Add this guard clause:
	if user.Status != "active" {
		return errors.New("account not confirmed")
	}

	var companyID *int

	// Only allow setting companyID if user is internal
	if strings.HasSuffix(strings.ToLower(user.Email), "@madagascarairlines.com") {
		if req.CompanyID != nil {
			companyID = req.CompanyID
		}
	}

	return s.repo.UpdateUserProfile(ctx, userID, req.FullName, companyID)
}
