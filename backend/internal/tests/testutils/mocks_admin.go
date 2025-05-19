// Package testutils provides mock implementations for unit and integration testing.
package testutils

import (
	"context"

	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepo mocks the admin repository interface
type MockAdminRepo struct {
	mock.Mock
}

// CreateUser mocks the creation of a user record.
func (m *MockAdminRepo) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// IsEmailExists mocks the check for an existing email in the database.
func (m *MockAdminRepo) IsEmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

// FindUserIDByEmail mocks fetching a user's ID by their email address.
func (m *MockAdminRepo) FindUserIDByEmail(ctx context.Context, email string) (int64, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(int64), args.Error(1)
}

// SetConfirmationToken mocks the setting of a confirmation token for a user.
func (m *MockAdminRepo) SetConfirmationToken(ctx context.Context, userID int64, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

// MockHasher is a mock hasher for testing
type MockHasher struct {
	mock.Mock
}

// HashPassword mocks the password hashing function.
func (m *MockHasher) HashPassword(p string) (string, error) {
	args := m.Called(p)
	return args.String(0), args.Error(1)
}
