package admin

import (
	"context"
	"errors"
	"testing"

	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAdminRepo is a mock implementation of the AdminRepo interface.
type MockAdminRepo struct {
	mock.Mock
}

// MockHasher is a mock implementation of the Hasher interface.
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAdminRepo) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	mockHasher := new(MockHasher)
	service := NewAdminService(mockRepo, mockHasher)

	req := CreateUserRequest{
		Email:    "user@madagascarairlines.com",
		Password: "securepass",
	}

	hashed := "hashedpassword"

	mockHasher.On("HashPassword", req.Password).Return(hashed, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
		return u.Email == req.Email && u.PasswordHash == hashed
	})).Return(nil)

	err := service.CreateUser(context.Background(), req, "admin")
	assert.NoError(t, err)

	mockHasher.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_HashFailure(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	mockHasher := new(MockHasher)
	service := NewAdminService(mockRepo, mockHasher)

	req := CreateUserRequest{
		Email:    "user@madagascarairlines.com",
		Password: "securepass",
	}

	mockHasher.On("HashPassword", req.Password).Return("", errors.New("hash error"))

	err := service.CreateUser(context.Background(), req, "admin")
	assert.EqualError(t, err, "hash error")

	mockHasher.AssertExpectations(t)
	mockRepo.AssertExpectations(t) // still called, though may be unused
}

func TestCreateUser_DBInsertFailure(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	mockHasher := new(MockHasher)
	service := NewAdminService(mockRepo, mockHasher)

	req := CreateUserRequest{
		Email:    "user@madagascarairlines.com",
		Password: "securepass",
	}

	hashed := "hashedpassword"

	mockHasher.On("HashPassword", req.Password).Return(hashed, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).
		Return(errors.New("insert failed"))

	err := service.CreateUser(context.Background(), req, "admin")
	assert.EqualError(t, err, "insert failed")

	mockHasher.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
