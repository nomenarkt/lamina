// internal/admin/mocks_test.go
package admin

import (
	"context"

	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/mock"
)

type MockAdminRepo struct {
	mock.Mock
}

func (m *MockAdminRepo) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(p string) (string, error) {
	args := m.Called(p)
	return args.String(0), args.Error(1)
}
