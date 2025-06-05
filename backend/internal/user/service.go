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

// UpdateUserProfile updates a user's full name and optionally their employee ID,
// depending on whether the user's email belongs to the internal domain.
func (s *Service) UpdateUserProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Status != "active" {
		return errors.New("account not confirmed")
	}

	// Only allow setting employee ID if user is internal
	if !strings.HasSuffix(strings.ToLower(user.Email), "@madagascarairlines.com") {
		req.EmployeeID = nil
	}

	return s.repo.UpdateUserProfile(ctx, userID, req.FullName, req.EmployeeID, req.Phone, req.Address)
}

// CompleteProfileByUserType validates profile fields based on user type and updates user.
func (s *Service) CompleteProfileByUserType(ctx context.Context, userID int64, req UpdateProfileRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Status != "active" {
		return errors.New("account not confirmed")
	}

	switch user.UserType {
	case "external":
		if req.FullName == "" {
			return errors.New("external users must provide name")
		}
	case "internal":
		if req.FullName == "" || req.EmployeeID == nil || req.Phone == nil || req.Address == nil {
			return errors.New("internal users must provide full name, employee ID, phone, and address")
		}
	}

	return s.repo.UpdateUserProfile(ctx, userID, req.FullName, req.EmployeeID, req.Phone, req.Address)
}

// IsAdmin checks if a user is an admin.
func (s *Service) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	return s.repo.IsAdmin(ctx, userID)
}

// MarkUserActive sets the user status to active.
func (s *Service) MarkUserActive(ctx context.Context, userID int64) error {
	return s.repo.MarkUserActive(ctx, userID)
}

// DeleteExpiredPendingUsers removes users with pending status that have expired.
func (s *Service) DeleteExpiredPendingUsers(ctx context.Context) error {
	return s.repo.DeleteExpiredPendingUsers(ctx)
}

// CreateUser creates a new user in the system and stores it in the repository.
func (s *Service) CreateUser(ctx context.Context, user *User) error {
	existing, err := s.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already in use")
	}
	return s.repo.Create(ctx, user)
}
