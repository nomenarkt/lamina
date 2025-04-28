package auth

import (
	"context"
	"errors"

	"github.com/nomenarkt/lamina/common/utils"
)

type AuthService struct {
	repo *AuthRepository
}

func NewAuthService(r *AuthRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Signup(ctx context.Context, req SignupRequest) (AuthResponse, error) {
	existingUser, _ := s.repo.FindByEmail(ctx, req.Email)
	if existingUser.ID != 0 {
		return AuthResponse{}, errors.New("email already in use")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return AuthResponse{}, err
	}

	userID, err := s.repo.CreateUser(ctx, req.Email, hashedPassword)
	if err != nil {
		return AuthResponse{}, err
	}

	accessToken, refreshToken, err := utils.GenerateTokens(userID, req.Email)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return AuthResponse{}, errors.New("invalid email or password")
	}

	if err := utils.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		return AuthResponse{}, errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email)
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
