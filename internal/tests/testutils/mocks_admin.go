// Package testutils provides reusable mock implementations for integration tests.
package testutils

import (
	"context"

	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepo is a mock implementation of the Admin repository used for testing.
type MockAdminRepo struct {
	mock.Mock
}

// CreateUser mocks the creation of a new user in the admin repo.
func (m *MockAdminRepo) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

// MockHasher is a mock implementation of a password hashing interface.
type MockHasher struct {
	mock.Mock
}

// HashPassword mocks the hashing of a password.
func (m *MockHasher) HashPassword(p string) (string, error) {
	args := m.Called(p)
	return args.String(0), args.Error(1)
}
