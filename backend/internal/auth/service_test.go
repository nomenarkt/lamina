package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) CreateUserInvite(ctx context.Context, email, userType string, accessExpires *time.Time) (int64, error) {
	args := m.Called(ctx, email, userType, accessExpires)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuthRepo) IsEmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepo) CreateUser(ctx context.Context, companyID int, email string, hash string) (int64, error) {
	args := m.Called(ctx, companyID, email, hash)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuthRepo) CreateUserWithType(ctx context.Context, companyID *int, email, hash, userType string) (int64, error) {
	args := m.Called(ctx, companyID, email, hash, userType)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuthRepo) FindByEmail(ctx context.Context, email string) (user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockAuthRepo) FindByConfirmationToken(ctx context.Context, token string) (user.User, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockAuthRepo) MarkUserConfirmed(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAuthRepo) SetConfirmationToken(ctx context.Context, userID int64, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockAuthRepo) UpdatePasswordAndActivate(ctx context.Context, userID int64, hashed string) error {
	args := m.Called(ctx, userID, hashed)
	return args.Error(0)
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockAuthRepo)
	u := user.User{ID: 1, Email: "test@example.com", PasswordHash: "any", Status: "active"}
	repo.On("FindByEmail", mock.Anything, "test@example.com").Return(u, nil)

	service := &Service{
		repo: repo,
		checkPassword: func(_, _ string) error {
			fmt.Println("✅ checkPassword mock called")
			return nil
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "access-token", "refresh-token", nil
		},
	}

	resp, err := service.Login(context.Background(), LoginRequest{
		Email:    "test@example.com",
		Password: "whatever",
	})

	assert.NoError(t, err)
	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := new(MockAuthRepo)
	u := user.User{ID: 1, Email: "test@example.com", PasswordHash: "wrong", Status: "active"}
	repo.On("FindByEmail", mock.Anything, "test@example.com").Return(u, nil)

	service := &Service{
		repo: repo,
		checkPassword: func(_, _ string) error {
			return errors.New("invalid")
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "", "", nil
		},
	}

	_, err := service.Login(context.Background(), LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpass",
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid email or password", err.Error())
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := new(MockAuthRepo)
	repo.On("FindByEmail", mock.Anything, "missing@example.com").Return(user.User{}, errors.New("not found"))

	service := &Service{
		repo:          repo,
		checkPassword: func(_, _ string) error { return nil },
		generateTokens: func(_ user.User) (string, string, error) {
			return "token", "refresh", nil
		},
	}

	_, err := service.Login(context.Background(), LoginRequest{
		Email:    "missing@example.com",
		Password: "irrelevant",
	})

	assert.Error(t, err)
	assert.Equal(t, "invalid email or password", err.Error())
}

func TestSignupUser_Success(t *testing.T) {
	repo := new(MockAuthRepo)
	email := "user@madagascarairlines.com"

	repo.On("IsEmailExists", email).Return(false, nil)
	repo.On("CreateUserWithType", mock.Anything, (*int)(nil), email, "hashed123", "internal").Return(int64(42), nil)
	repo.On("SetConfirmationToken", mock.Anything, int64(42), mock.AnythingOfType("string")).Return(nil)

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return "hashed123", nil
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "access-token", "refresh-token", nil
		},
	}

	resp, err := service.SignupUser(context.Background(), SignupRequest{
		Email:    email,
		Password: "mypassword",
	})

	assert.NoError(t, err)
	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)
}

func TestSignupUser_EmailExists(t *testing.T) {
	repo := new(MockAuthRepo)
	email := "user@madagascarairlines.com"

	repo.On("IsEmailExists", email).Return(true, nil)

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return "ignored", nil
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "", "", nil
		},
	}

	_, err := service.SignupUser(context.Background(), SignupRequest{
		Email:    email,
		Password: "irrelevant",
	})

	assert.Error(t, err)
	assert.Equal(t, "email already registered", err.Error())
}

func TestSignupUser_HashFailure(t *testing.T) {
	repo := new(MockAuthRepo)
	email := "user@madagascarairlines.com"
	repo.On("IsEmailExists", email).Return(false, nil)

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return "", errors.New("hash failed")
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "", "", nil
		},
	}

	_, err := service.SignupUser(context.Background(), SignupRequest{
		Email:    email,
		Password: "bad",
	})

	assert.Error(t, err)
	assert.EqualError(t, err, "failed to hash password: hash failed")
}

func TestSignupUser_TokenFailure(t *testing.T) {
	repo := new(MockAuthRepo)
	email := "user@madagascarairlines.com"
	repo.On("IsEmailExists", email).Return(false, nil)
	repo.On("CreateUserWithType", mock.Anything, (*int)(nil), email, "hashedok", "internal").Return(int64(99), nil)
	repo.On("SetConfirmationToken", mock.Anything, int64(99), mock.AnythingOfType("string")).Return(errors.New("mock token failure"))

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return "hashedok", nil
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "", "", errors.New("token failed")
		},
	}

	_, err := service.SignupUser(context.Background(), SignupRequest{
		Email:    email,
		Password: "good",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to issue token: mock token failure")
}

func TestSignupUser_InvalidDomain(t *testing.T) {
	repo := new(MockAuthRepo)
	email := "intruder@notallowed.com"

	repo.On("IsEmailExists", email).Return(false, nil)

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return "irrelevant", nil
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "", "", nil
		},
	}

	_, err := service.SignupUser(context.Background(), SignupRequest{
		Email:    email,
		Password: "shouldfail",
	})

	assert.Error(t, err)
	assert.Equal(t, "only @madagascarairlines.com emails are allowed", err.Error())
}

func TestCompleteInvite_Success(t *testing.T) {
	repo := new(MockAuthRepo)
	userID := int64(77)
	token := "valid-token"
	password := "newpass"
	hashed := "hashedpass"

	testUser := user.User{
		ID:        userID,
		Email:     "invitee@example.com",
		Status:    "pending",
		CreatedAt: time.Now(), // 👈 prevents TTL expiry
	}

	repo.On("FindByConfirmationToken", mock.Anything, token).Return(testUser, nil)
	repo.On("UpdatePasswordAndActivate", mock.Anything, userID, hashed).Return(nil)

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return hashed, nil
		},
		generateTokens: func(_ user.User) (string, string, error) {
			return "access-token", "refresh-token", nil
		},
		confirmationTTL: 24 * time.Hour,
	}

	resp, err := service.CompleteInvite(context.Background(), token, password)
	assert.NoError(t, err)
	assert.Equal(t, "access-token", resp.AccessToken)
	assert.Equal(t, "refresh-token", resp.RefreshToken)
}

func TestConfirmRegistration_Success(t *testing.T) {
	repo := new(MockAuthRepo)
	repo.On("FindByConfirmationToken", mock.Anything, "valid-token").Return(user.User{
		ID:        1,
		Status:    "pending",
		CreatedAt: time.Now(), // 👈 within TTL
	}, nil)
	repo.On("MarkUserConfirmed", mock.Anything, int64(1)).Return(nil)

	service := &Service{
		repo:            repo,
		confirmationTTL: 24 * time.Hour,
	}

	err := service.ConfirmRegistration(context.Background(), "valid-token")
	assert.NoError(t, err)
}
