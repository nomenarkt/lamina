package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepo implements UserRepository interface for testing
type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByID(ctx context.Context, id int64) (User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(User), args.Error(1)
}

func (m *MockUserRepo) FindAll(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func TestGetMe_Success(t *testing.T) {
	repo := new(MockUserRepo)
	service := NewUserService(repo)

	expectedUser := User{ID: 1, Email: "me@example.com", Role: "user"}
	repo.On("FindByID", mock.Anything, int64(1)).Return(expectedUser, nil)

	user, err := service.GetMe(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Role, user.Role)
}

func TestGetMe_UserNotFound(t *testing.T) {
	repo := new(MockUserRepo)
	service := NewUserService(repo)

	repo.On("FindByID", mock.Anything, int64(42)).Return(User{}, errors.New("not found"))

	_, err := service.GetMe(context.Background(), 42)

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}
