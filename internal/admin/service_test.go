package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/assert"
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

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	repo := new(MockAdminRepo)
	hasher := new(MockHasher)
	service := NewAdminService(repo, hasher)

	req := CreateUserRequest{
		Email:    "user@madagascarairlines.com",
		Password: "securepass",
		Role:     "admin",
	}

	hashedPassword := "hashedpassword"

	hasher.On("HashPassword", req.Password).Return(hashedPassword, nil)
	repo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
		return u.Email == "user@madagascarairlines.com" && u.PasswordHash == "hashedpassword" && u.Role == "admin"
	})).Return(nil)

	err := service.CreateUser(context.Background(), req, "admin")
	assert.NoError(t, err)
}

func TestCreateUser_HashFailure(t *testing.T) {
	repo := new(MockAdminRepo)
	hasher := new(MockHasher)
	service := NewAdminService(repo, hasher)

	req := CreateUserRequest{
		Email:    "user@madagascarairlines.com",
		Password: "securepass",
		Role:     "admin",
	}

	hasher.On("HashPassword", req.Password).Return("", errors.New("hash error"))

	err := service.CreateUser(context.Background(), req, "admin")
	assert.EqualError(t, err, "hash error")
}

func TestCreateUser_DBInsertFailure(t *testing.T) {
	repo := new(MockAdminRepo)
	hasher := new(MockHasher)
	service := NewAdminService(repo, hasher)

	req := CreateUserRequest{
		Email:    "user@madagascarairlines.com",
		Password: "securepass",
		Role:     "admin",
	}

	hashedPassword := "hashedpassword"

	hasher.On("HashPassword", req.Password).Return(hashedPassword, nil)
	repo.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("insert failed"))

	err := service.CreateUser(context.Background(), req, "admin")
	assert.EqualError(t, err, "insert failed")
}
