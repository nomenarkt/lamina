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

func (m *MockAdminRepo) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockAdminRepo) IsEmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockAdminRepo) FindUserIDByEmail(ctx context.Context, email string) (int64, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAdminRepo) SetConfirmationToken(ctx context.Context, userID int64, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

// MockHasher is a mock implementation of the Hasher interface.
type MockHasher struct {
	mock.Mock
}

func (m *MockHasher) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func TestInviteUser_Success(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	mockHasher := new(MockHasher)
	service := NewAdminService(mockRepo, mockHasher)

	req := CreateUserRequest{Email: "newuser@madagascarairlines.com"}

	mockRepo.On("IsEmailExists", req.Email).Return(false, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("FindUserIDByEmail", mock.Anything, req.Email).Return(int64(99), nil)
	mockRepo.On("SetConfirmationToken", mock.Anything, int64(99), mock.Anything).Return(nil)

	err := service.InviteUser(context.Background(), req, "admin")
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestInviteUser_EmailExists(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	mockHasher := new(MockHasher)
	service := NewAdminService(mockRepo, mockHasher)

	req := CreateUserRequest{Email: "existing@madagascarairlines.com"}

	mockRepo.On("IsEmailExists", req.Email).Return(true, nil)

	err := service.InviteUser(context.Background(), req, "admin")
	assert.EqualError(t, err, "email already registered")
}

func TestInviteUser_RepoFails(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	mockHasher := new(MockHasher)
	service := NewAdminService(mockRepo, mockHasher)

	req := CreateUserRequest{Email: "dbfail@madagascarairlines.com"}

	mockRepo.On("IsEmailExists", req.Email).Return(false, nil)
	mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("DB insert error"))

	err := service.InviteUser(context.Background(), req, "admin")
	assert.EqualError(t, err, "DB insert error")
}
