package user

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	user := args.Get(0).(User)
	return &user, args.Error(1)
}

func (m *MockUserRepo) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id int64) (*User, error) {
	args := m.Called(id)
	user := args.Get(0).(User)
	return &user, args.Error(1)
}

func (m *MockUserRepo) FindAll(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockUserRepo) UpdateFullName(ctx context.Context, userID int64, fullName string) error {
	args := m.Called(ctx, userID, fullName)
	return args.Error(0)
}

func (m *MockUserRepo) UpdateUserProfile(ctx context.Context, userID int64, fullName string, companyID *int) error {
	args := m.Called(ctx, userID, fullName, companyID)
	return args.Error(0)
}

func TestGetMe_Success(t *testing.T) {
	repo := new(MockUserRepo)
	service := NewUserService(repo)

	expectedUser := User{ID: 3190, Email: "user@madagascarairlines.com"}
	repo.On("FindByID", int64(3190)).Return(expectedUser, nil)

	user, err := service.GetMe(context.Background(), 3190)
	assert.NoError(t, err)
	assert.Equal(t, user, &expectedUser)
}

func TestGetMe_InvalidID(t *testing.T) {
	repo := new(MockUserRepo)
	service := NewUserService(repo)

	repo.On("FindByID", int64(0)).Return(User{}, errors.New("user not found"))

	_, err := service.GetMe(context.Background(), 0)
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestFindAll_UsersExist(t *testing.T) {
	repo := new(MockUserRepo)
	service := NewUserService(repo)

	mockUsers := []User{
		{ID: 3190, Email: "first@madagascarairlines.com"},
		{ID: 3191, Email: "second@madagascarairlines.com"},
	}

	repo.On("FindAll", mock.Anything).Return(mockUsers, nil)

	users, err := service.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "first@madagascarairlines.com", users[0].Email)
}

func TestFindAll_NoUsers(t *testing.T) {
	repo := new(MockUserRepo)
	service := NewUserService(repo)

	repo.On("FindAll", mock.Anything).Return([]User{}, nil)

	users, err := service.FindAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, users)
}
