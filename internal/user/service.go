package user

import (
	"context"
	"strings"
)

type UserService struct {
	repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetMe(ctx context.Context, id int64) (*User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) FindAll(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.FindAll(ctx)
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID int64, req UpdateProfileRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	var companyID *int

	// If the user has a corporate email, allow setting company_id (if provided)
	if strings.HasSuffix(strings.ToLower(user.Email), "@madagascarairlines.com") {
		if req.CompanyID != nil {
			companyID = req.CompanyID
		}
	}

	return s.repo.UpdateUserProfile(ctx, userID, req.FullName, companyID)
}
