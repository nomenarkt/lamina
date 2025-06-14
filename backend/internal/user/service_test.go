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
	return args.Get(0).(*User), args.Error(1) // ✅ fixed
}

func (m *MockUserRepo) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, id int64) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	user := args.Get(0).(User)
	return &user, args.Error(1)
}

func (m *MockUserRepo) FindAll(ctx context.Context) ([]User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockUserRepo) UpdateUserProfile(ctx context.Context, userID int64, fullName string, employeeID *int, phone, address *string) error {
	args := m.Called(ctx, userID, fullName, employeeID, phone, address)
	return args.Error(0)
}

func (m *MockUserRepo) MarkUserActive(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteExpiredPendingUsers(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func newMockedUserService() (*Service, *MockUserRepo) {
	repo := new(MockUserRepo)
	svc := NewUserService(repo)
	return svc, repo
}

func TestUserService_GetMe(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo := newMockedUserService()

		expected := User{ID: 3190, Email: "user@madagascarairlines.com"}
		repo.On("FindByID", context.Background(), int64(3190)).Return(expected, nil)

		user, err := svc.GetMe(context.Background(), 3190)
		assert.NoError(t, err)
		assert.Equal(t, &expected, user)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo := newMockedUserService()

		repo.On("FindByID", context.Background(), int64(0)).Return(User{}, errors.New("user not found"))

		_, err := svc.GetMe(context.Background(), 0)
		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestUserService_FindAll(t *testing.T) {
	t.Run("users exist", func(t *testing.T) {
		svc, repo := newMockedUserService()

		mockUsers := []User{
			{ID: 3190, Email: "first@madagascarairlines.com"},
			{ID: 3191, Email: "second@madagascarairlines.com"},
		}
		repo.On("FindAll", context.Background()).Return(mockUsers, nil)

		users, err := svc.FindAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "first@madagascarairlines.com", users[0].Email)
	})

	t.Run("no users", func(t *testing.T) {
		svc, repo := newMockedUserService()

		repo.On("FindAll", context.Background()).Return([]User{}, nil)

		users, err := svc.FindAll(context.Background())
		assert.NoError(t, err)
		assert.Empty(t, users)
	})
}

func TestUserService_CompleteProfileByUserType(t *testing.T) {
	type testCase struct {
		name       string
		user       User
		request    UpdateProfileRequest
		expectErr  string
		expectCall bool
	}

	phone := "123456"
	address := "HQ"
	employeeID := 1234

	tests := []testCase{
		{
			name: "external valid",
			user: User{ID: 1, Status: "active", UserType: "external"},
			request: UpdateProfileRequest{
				FullName: "John Doe",
			},
			expectErr:  "",
			expectCall: true,
		},
		{
			name: "internal valid",
			user: User{ID: 2, Status: "active", UserType: "internal"},
			request: UpdateProfileRequest{
				FullName:   "Jane Smith",
				EmployeeID: &employeeID,
				Phone:      &phone,
				Address:    &address,
			},
			expectErr:  "",
			expectCall: true,
		},
		{
			name: "internal missing phone",
			user: User{ID: 3, Status: "active", UserType: "internal"},
			request: UpdateProfileRequest{
				FullName:   "Sam MissingPhone",
				EmployeeID: &employeeID,
				Address:    &address,
			},
			expectErr:  "internal users must provide full name, employee ID, phone, and address",
			expectCall: false,
		},
		{
			name: "external missing name",
			user: User{ID: 4, Status: "active", UserType: "external"},
			request: UpdateProfileRequest{
				FullName: "",
			},
			expectErr:  "external users must provide name",
			expectCall: false,
		},
		{
			name: "not confirmed user",
			user: User{ID: 5, Status: "pending", UserType: "external"},
			request: UpdateProfileRequest{
				FullName: "Ghost User",
			},
			expectErr:  "account not confirmed",
			expectCall: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, repo := newMockedUserService()

			repo.On("FindByID", mock.Anything, tc.user.ID).Return(tc.user, nil)

			if tc.expectCall {
				repo.On("UpdateUserProfile", mock.Anything, tc.user.ID, tc.request.FullName, tc.request.EmployeeID, tc.request.Phone, tc.request.Address).Return(nil)
			}

			err := svc.CompleteProfileByUserType(context.Background(), tc.user.ID, tc.request)

			if tc.expectErr != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_IsAdmin(t *testing.T) {
	svc, repo := newMockedUserService()

	repo.On("IsAdmin", context.Background(), int64(77)).Return(true, nil)

	isAdmin, err := svc.IsAdmin(context.Background(), 77)
	assert.NoError(t, err)
	assert.True(t, isAdmin)
}

func TestUserService_MarkUserActive(t *testing.T) {
	svc, repo := newMockedUserService()

	repo.On("MarkUserActive", context.Background(), int64(2024)).Return(nil)

	err := svc.MarkUserActive(context.Background(), 2024)
	assert.NoError(t, err)
}

func TestUserService_DeleteExpiredPendingUsers(t *testing.T) {
	svc, repo := newMockedUserService()

	repo.On("DeleteExpiredPendingUsers", context.Background()).Return(nil)

	err := svc.DeleteExpiredPendingUsers(context.Background())
	assert.NoError(t, err)
}

func (m *MockUserRepo) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc, repo := newMockedUserService()

		input := &User{
			Email:        "new@user.com",
			PasswordHash: "hashedpass",
			Role:         "user",
			Status:       "pending",
		}

		repo.On("FindByEmail", mock.Anything, input.Email).Return(nil, nil)
		repo.On("Create", mock.Anything, input).Return(nil)

		err := svc.CreateUser(context.Background(), input)
		assert.NoError(t, err)
	})

	t.Run("duplicate email", func(t *testing.T) {
		svc, repo := newMockedUserService()

		existing := User{ID: 999, Email: "dupe@user.com"}
		input := &User{Email: "dupe@user.com"}

		repo.On("FindByEmail", mock.Anything, input.Email).Return(&existing, nil)

		err := svc.CreateUser(context.Background(), input)
		assert.Error(t, err)
		assert.Equal(t, "email already in use", err.Error())
	})
}
