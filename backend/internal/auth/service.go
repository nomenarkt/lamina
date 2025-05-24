package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/user"
	"golang.org/x/crypto/bcrypt"
)

// InviteRequest represents a request to invite a user with optional duration for access.
type InviteRequest struct {
	Email    string `json:"email" binding:"required,email"`
	UserType string `json:"user_type" binding:"required"` // internal or external
	Duration string `json:"duration,omitempty"`           // only for external (e.g., "1w", "90m")
}

// ServiceInterface defines the interface for Service business logic.
type ServiceInterface interface {
	Signup(c *gin.Context)
	InviteUser(c *gin.Context)
	SignupUser(ctx context.Context, req SignupRequest) (Response, error)
	Login(ctx context.Context, req LoginRequest) (Response, error)
	CompleteInvite(ctx context.Context, token string, password string) (Response, error)
	ConfirmRegistration(ctx context.Context, token string) error
}

// Service provides user authentication logic.
type Service struct {
	repo           Repository
	checkPassword  func(raw, hash string) error
	hashPassword   func(p string) (string, error)
	generateTokens func(u user.User) (string, string, error)
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

// InviteUser handles inviting a user (internal or external) via email and optional duration.
func (s *Service) InviteUser(c *gin.Context) {
	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.UserType != "internal" && req.UserType != "external" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_type"})
		return
	}

	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	var accessExpires *time.Time
	if req.UserType == "external" && req.Duration != "" {
		duration, err := ParseFlexibleDuration(req.Duration)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format"})
			return
		}
		t := time.Now().Add(duration)
		accessExpires = &t
	}

	userID, err := s.repo.CreateUserInvite(c.Request.Context(), req.Email, req.UserType, accessExpires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}
	if err := s.repo.SetConfirmationToken(c.Request.Context(), userID, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store token"})
		return
	}
	if err := SendConfirmationEmail(req.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invite email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User invited"})
}

// SignupUser registers a new user and returns tokens.
func (s *Service) SignupUser(ctx context.Context, req SignupRequest) (Response, error) {
	log.Printf("➡️ SignupUser called with email: %s", req.Email)

	if !strings.HasSuffix(req.Email, "@madagascarairlines.com") {
		return Response{}, errors.New("only @madagascarairlines.com emails are allowed")
	}

	log.Println("🔍 Checking if email exists...")
	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return Response{}, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return Response{}, errors.New("email already registered")
	}

	log.Println("🔐 Hashing password...")
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return Response{}, fmt.Errorf("failed to hash password: %w", err)
	}

	log.Println("📝 Creating user with type 'internal'...")
	var unsetCompanyID *int
	userID, err := s.repo.CreateUserWithType(ctx, unsetCompanyID, req.Email, hashedPassword, "internal")

	if err != nil {
		log.Printf("❌ DB error during CreateUserWithType: %v", err)
		return Response{}, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("📨 Generating confirmation token for user ID: %d", userID)
	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return Response{}, fmt.Errorf("failed to generate confirmation token: %w", err)
	}

	log.Println("💾 Storing confirmation token...")
	if err := s.repo.SetConfirmationToken(ctx, userID, token); err != nil {
		return Response{}, fmt.Errorf("failed to store confirmation token: %w", err)
	}

	log.Println("📧 Sending confirmation email...")
	if err := SendConfirmationEmail(req.Email, token); err != nil {
		return Response{}, fmt.Errorf("failed to send confirmation email: %w", err)
	}

	log.Println("🔑 Generating access and refresh tokens...")
	newUser := user.User{
		ID:    userID,
		Email: req.Email,
		Role:  "user",
	}
	access, refresh, err := s.generateTokens(newUser)
	if err != nil {
		return Response{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	log.Println("✅ Signup flow completed successfully")
	return Response{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// Login validates a user and returns new tokens.
func (s *Service) Login(ctx context.Context, req LoginRequest) (Response, error) {
	userRecord, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return Response{}, errors.New("invalid email or password")
	}

	if userRecord.Status != "active" {
		return Response{}, errors.New("account not confirmed")
	}

	if err := s.checkPassword(req.Password, userRecord.PasswordHash); err != nil {
		return Response{}, errors.New("invalid email or password")
	}

	access, refresh, err := s.generateTokens(userRecord)
	if err != nil {
		return Response{}, err
	}

	return Response{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// CompleteInvite sets password, activates user, and returns JWT tokens.
func (s *Service) CompleteInvite(ctx context.Context, token string, password string) (Response, error) {
	userRecord, err := s.repo.FindByConfirmationToken(ctx, token)
	if err != nil {
		return Response{}, errors.New("invalid or expired token")
	}
	if userRecord.Status != "pending" {
		return Response{}, errors.New("user already confirmed")
	}
	if time.Since(userRecord.CreatedAt) > 24*time.Hour {
		return Response{}, errors.New("token expired")
	}

	hashed, err := s.hashPassword(password)
	if err != nil {
		return Response{}, errors.New("failed to hash password")
	}

	if err := s.repo.UpdatePasswordAndActivate(ctx, userRecord.ID, hashed); err != nil {
		return Response{}, errors.New("failed to activate account")
	}

	access, refresh, err := s.generateTokens(userRecord)
	if err != nil {
		return Response{}, errors.New("failed to generate tokens")
	}

	return Response{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// ConfirmRegistration activates a pending user account from email token.
func (s *Service) ConfirmRegistration(ctx context.Context, token string) error {
	userRecord, err := s.repo.FindByConfirmationToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if userRecord.Status != "pending" {
		return errors.New("user already confirmed")
	}

	if time.Since(userRecord.CreatedAt) > 24*time.Hour {
		return errors.New("token expired")
	}

	if err := s.repo.MarkUserConfirmed(ctx, userRecord.ID); err != nil {
		return errors.New("failed to activate account")
	}

	return nil
}

// ValidateConfirmationToken validates and checks expiration of the confirmation token.
func (s *Service) ValidateConfirmationToken(ctx context.Context, token string) error {
	userRecord, err := s.repo.FindByConfirmationToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired confirmation token")
	}
	if userRecord.Status != "pending" {
		return errors.New("user already confirmed")
	}
	if time.Since(userRecord.CreatedAt) > 24*time.Hour {
		return errors.New("confirmation token expired")
	}
	return nil // No-op (user will confirm later via /complete-invite)
}

// HashPassword hashes a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ParseFlexibleDuration parses duration strings like "15m", "2h", "3d", "1w", "6mo", "1y".
func ParseFlexibleDuration(input string) (time.Duration, error) {
	switch {
	case strings.HasSuffix(input, "y"):
		n, err := strconv.Atoi(strings.TrimSuffix(input, "y"))
		return time.Hour * 24 * 365 * time.Duration(n), err
	case strings.HasSuffix(input, "mo"):
		n, err := strconv.Atoi(strings.TrimSuffix(input, "mo"))
		return time.Hour * 24 * 30 * time.Duration(n), err
	case strings.HasSuffix(input, "w"):
		n, err := strconv.Atoi(strings.TrimSuffix(input, "w"))
		return time.Hour * 24 * 7 * time.Duration(n), err
	case strings.HasSuffix(input, "d"):
		n, err := strconv.Atoi(strings.TrimSuffix(input, "d"))
		return time.Hour * 24 * time.Duration(n), err
	default:
		return time.ParseDuration(input) // handles "15m", "3h"
	}
}

// CheckPasswordHash compares a plaintext password to a bcrypt hash.
func CheckPasswordHash(_raw, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(_raw))
}
