package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/nomenarkt/lamina/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) IsEmailExists(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthRepo) CreateUser(ctx context.Context, companyID int, email string, hash string) (int64, error) {
	args := m.Called(ctx, companyID, email, hash)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockAuthRepo) FindByEmail(ctx context.Context, email string) (user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(user.User), args.Error(1)
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockAuthRepo)
	u := user.User{ID: 1, Email: "test@example.com", PasswordHash: "any"}
	repo.On("FindByEmail", mock.Anything, "test@example.com").Return(u, nil)

	service := &Service{
		repo: repo,
		checkPassword: func(_, _ string) error {
			fmt.Println("✅ checkPassword mock called")
			return nil
		},
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
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
	u := user.User{ID: 1, Email: "test@example.com", PasswordHash: "wrong"}
	repo.On("FindByEmail", mock.Anything, "test@example.com").Return(u, nil)

	service := &Service{
		repo: repo,
		checkPassword: func(_, _ string) error {
			return errors.New("invalid")
		},
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
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
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
			return "", "", nil
		},
	}

	_, err := service.Login(context.Background(), LoginRequest{
		Email:    "missing@example.com",
		Password: "irrelevant",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSignupUser_Success(t *testing.T) {
	repo := new(MockAuthRepo)
	email := "user@madagascarairlines.com"

	repo.On("IsEmailExists", email).Return(false, nil)
	repo.On("CreateUser", mock.Anything, 0, email, "hashed123").Return(int64(42), nil)

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return "hashed123", nil
		},
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
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
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
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
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
			return "", "", nil
		},
	}

	_, err := service.SignupUser(context.Background(), SignupRequest{
		Email:    email,
		Password: "bad",
	})

	assert.Error(t, err)
	assert.Equal(t, "failed to hash password", err.Error())
}

func TestSignupUser_TokenFailure(t *testing.T) {
	repo := new(MockAuthRepo)
	email := "user@madagascarairlines.com"
	repo.On("IsEmailExists", email).Return(false, nil)
	repo.On("CreateUser", mock.Anything, 0, email, "hashedok").Return(int64(99), nil)

	service := &Service{
		repo: repo,
		hashPassword: func(_ string) (string, error) {
			return "hashedok", nil
		},
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
			return "", "", errors.New("token failed")
		},
	}

	_, err := service.SignupUser(context.Background(), SignupRequest{
		Email:    email,
		Password: "good",
	})

	assert.Error(t, err)
	assert.Equal(t, "failed to generate tokens", err.Error())
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
		generateTokens: func(_ int64, _, _ string) (string, string, error) {
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
