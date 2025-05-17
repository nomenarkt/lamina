package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/utils"
	"golang.org/x/crypto/bcrypt"
)

// ServiceInterface defines the interface for Service business logic.
type ServiceInterface interface {
	Signup(c *gin.Context)
	SignupUser(ctx context.Context, req SignupRequest) (Response, error)
	Login(ctx context.Context, req LoginRequest) (Response, error)
}

// Service provides user authentication logic.
type Service struct {
	repo           Repository
	checkPassword  func(raw, hash string) error
	hashPassword   func(p string) (string, error)
	generateTokens func(id int64, email, role string) (string, string, error)
}

// NewService creates a new Service instance.
func NewService(r Repository) *Service {
	return &Service{
		repo:           r,
		checkPassword:  CheckPasswordHash,
		hashPassword:   HashPassword,
		generateTokens: GenerateTokensFromEnv,
	}
}

// SignupUser registers a new user and returns tokens.
func (s *Service) SignupUser(ctx context.Context, req SignupRequest) (Response, error) {
	if !strings.HasSuffix(req.Email, "@madagascarairlines.com") {
		return Response{}, errors.New("only @madagascarairlines.com emails are allowed")
	}

	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return Response{}, errors.New("failed to check user existence")
	}
	if exists {
		return Response{}, errors.New("email already registered")
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return Response{}, errors.New("failed to hash password")
	}

	userID, err := s.repo.CreateUser(ctx, 0, req.Email, hashedPassword)
	if err != nil {
		return Response{}, errors.New("failed to create user")
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return Response{}, errors.New("failed to generate confirmation token")
	}
	if err := s.repo.SetConfirmationToken(ctx, userID, token); err != nil {
		return Response{}, errors.New("failed to store confirmation token")
	}
	if err := SendConfirmationEmail(req.Email, token); err != nil {
		return Response{}, errors.New("failed to send confirmation email")
	}

	access, refresh, err := s.generateTokens(userID, req.Email, "user")
	if err != nil {
		return Response{}, errors.New("failed to generate tokens")
	}

	return Response{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// Signup is a Gin handler that registers a new user.
func (s *Service) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if !strings.HasSuffix(strings.ToLower(req.Email), "@madagascarairlines.com") {
		c.JSON(http.StatusForbidden, gin.H{"error": "signup is restricted to company email addresses"})
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

// Login validates a user and returns new tokens.
func (s *Service) Login(ctx context.Context, req LoginRequest) (Response, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return Response{}, errors.New("invalid email or password")
	}

	if user.Status != "active" {
		return Response{}, errors.New("account not confirmed")
	}

	if err := s.checkPassword(req.Password, user.PasswordHash); err != nil {
		return Response{}, errors.New("invalid email or password")
	}

	access, refresh, err := s.generateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		return Response{}, err
	}

	return Response{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password to a bcrypt hash.
func CheckPasswordHash(raw, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(raw))
}

// ConfirmRegistration validates and finalizes email confirmation via token.
func (s *Service) ConfirmRegistration(ctx context.Context, token string) error {
	user, err := s.repo.FindByConfirmationToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired confirmation token")
	}
	if user.Status != "pending" {
		return errors.New("user already confirmed")
	}
	if time.Since(user.CreatedAt) > 24*time.Hour {
		return errors.New("confirmation token expired")
	}
	return s.repo.MarkUserConfirmed(ctx, user.ID)
}
