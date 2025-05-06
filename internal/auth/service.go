package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/user"
)

type AuthRepoInterface interface {
	IsEmailExists(email string) (bool, error)
	CreateUser(ctx context.Context, companyID int, email string, hash string) (int64, error)
	FindByEmail(ctx context.Context, email string) (user.User, error)
}

type AuthService struct {
	repo           AuthRepoInterface
	checkPassword  func(raw, hash string) error
	hashPassword   func(p string) (string, error)
	generateTokens func(id int64, email, role string) (string, string, error)
}

func NewAuthService(r AuthRepoInterface) *AuthService {
	return &AuthService{
		repo:           r,
		checkPassword:  utils.CheckPasswordHash,
		hashPassword:   utils.HashPassword,
		generateTokens: utils.GenerateTokens,
	}
}

func (s *AuthService) SignupUser(ctx context.Context, req SignupRequest) (AuthResponse, error) {
	if !strings.HasSuffix(req.Email, "@madagascarairlines.com") {
		return AuthResponse{}, errors.New("only @madagascarairlines.com emails are allowed")
	}

	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return AuthResponse{}, errors.New("failed to check user existence")
	}
	if exists {
		return AuthResponse{}, errors.New("email already registered")
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return AuthResponse{}, errors.New("failed to hash password")
	}

	userID, err := s.repo.CreateUser(ctx, req.CompanyID, req.Email, hashedPassword)
	if err != nil {
		return AuthResponse{}, errors.New("failed to create user")
	}

	access, refresh, err := s.generateTokens(userID, req.Email, "user")
	if err != nil {
		return AuthResponse{}, errors.New("failed to generate tokens")
	}

	return AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *AuthService) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	resp, err := s.SignupUser(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
	})
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return AuthResponse{}, err
	}

	if err := s.checkPassword(req.Password, user.PasswordHash); err != nil {
		return AuthResponse{}, errors.New("invalid email or password")
	}

	access, refresh, err := s.generateTokens(user.ID, user.Email, user.Role) // ← use user.Role here
	if err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
