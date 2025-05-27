package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nomenarkt/lamina/common/utils"
	"github.com/nomenarkt/lamina/internal/user"
	"golang.org/x/crypto/bcrypt"
)

// InviteRequest represents a request payload for inviting a user (internal or external).
type InviteRequest struct {
	Email    string `json:"email" binding:"required,email"`
	UserType string `json:"user_type" binding:"required"` // internal or external
	Duration string `json:"duration,omitempty"`           // only for external (e.g., "1w", "90m")
}

// ServiceInterface defines all the methods available in the auth service layer.
type ServiceInterface interface {
	Signup(c *gin.Context)
	InviteUser(c *gin.Context)
	SignupUser(ctx context.Context, req SignupRequest) (Response, error)
	Login(ctx context.Context, req LoginRequest) (Response, error)
	CompleteInvite(ctx context.Context, token string, password string) (Response, error)
	ConfirmRegistration(ctx context.Context, token string) error
	ResendConfirmation(ctx context.Context, email string) error
}

// Service provides the implementation for authentication-related business logic.
type Service struct {
	repo            Repository
	checkPassword   func(raw, hash string) error
	hashPassword    func(p string) (string, error)
	generateTokens  func(u user.User) (string, string, error)
	confirmationTTL time.Duration
}

// NewService creates a new instance of the auth Service.
func NewService(r Repository) *Service {
	ttlHours := 24
	if envTTL := os.Getenv("CONFIRMATION_TOKEN_TTL_HOURS"); envTTL != "" {
		if parsed, err := strconv.Atoi(envTTL); err == nil {
			ttlHours = parsed
		}
	}
	return &Service{
		repo:            r,
		checkPassword:   CheckPasswordHash,
		hashPassword:    HashPassword,
		generateTokens:  GenerateTokensFromEnv,
		confirmationTTL: time.Duration(ttlHours) * time.Hour,
	}
}

// Signup handles registration of internal users using a corporate email.
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

// SignupUser creates a new internal user account and sends a confirmation email.
func (s *Service) SignupUser(ctx context.Context, req SignupRequest) (Response, error) {
	if !strings.HasSuffix(req.Email, "@madagascarairlines.com") {
		return Response{}, errors.New("only @madagascarairlines.com emails are allowed")
	}

	exists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return Response{}, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return Response{}, errors.New("email already registered")
	}

	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		return Response{}, fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := s.repo.CreateUserWithType(ctx, nil, req.Email, hashedPassword, "internal")
	if err != nil {
		return Response{}, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.issueConfirmationToken(ctx, userID)
	if err != nil {
		return Response{}, fmt.Errorf("failed to issue token: %w", err)
	}
	if err := s.notifyUserWithToken(req.Email, token, false); err != nil {
		return Response{}, fmt.Errorf("failed to send confirmation email: %w", err)
	}

	access, refresh, err := s.generateTokens(user.User{ID: userID, Email: req.Email, Role: "user"})
	if err != nil {
		return Response{}, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return Response{AccessToken: access, RefreshToken: refresh}, nil
}

// InviteUser allows admins to invite internal/external users with optional access duration.
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

	token, err := s.issueConfirmationToken(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}
	if err := s.notifyUserWithToken(req.Email, token, false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invite email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User invited"})
}

// ResendConfirmation resends the email confirmation link for a pending internal user.
func (s *Service) ResendConfirmation(ctx context.Context, email string) error {
	userRecord, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	log.Printf("DEBUG: Resend - email=%s user_type=%s status=%s", userRecord.Email, userRecord.UserType, userRecord.Status)

	if userRecord.Status != "pending" {
		return errors.New("user already confirmed or invalid status")
	}
	if userRecord.UserType != "internal" || !strings.HasSuffix(strings.ToLower(userRecord.Email), "@madagascarairlines.com") {
		return errors.New("resend allowed only for internal users")
	}

	token, err := s.issueConfirmationToken(ctx, userRecord.ID)
	if err != nil {
		return fmt.Errorf("failed to issue token: %w", err)
	}
	if err := s.notifyUserWithToken(userRecord.Email, token, true); err != nil {
		return fmt.Errorf("failed to send confirmation email: %w", err)
	}
	return nil
}

// CompleteInvite sets the password and activates an invited user.
func (s *Service) CompleteInvite(ctx context.Context, token string, password string) (Response, error) {
	userRecord, err := s.repo.FindByConfirmationToken(ctx, token)
	if err != nil {
		return Response{}, errors.New("invalid or expired token")
	}
	if userRecord.Status != "pending" {
		return Response{}, errors.New("user already confirmed")
	}
	if time.Since(userRecord.CreatedAt) > s.confirmationTTL {
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

	return Response{AccessToken: access, RefreshToken: refresh}, nil
}

// Login authenticates a user with email and password, returning tokens.
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
	return Response{AccessToken: access, RefreshToken: refresh}, nil
}

// ConfirmRegistration confirms a user account using a token sent via email.
func (s *Service) ConfirmRegistration(ctx context.Context, token string) error {
	userRecord, err := s.repo.FindByConfirmationToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}
	if userRecord.Status != "pending" {
		return errors.New("user already confirmed")
	}
	if time.Since(userRecord.CreatedAt) > s.confirmationTTL {
		return errors.New("token expired")
	}
	return s.repo.MarkUserConfirmed(ctx, userRecord.ID)
}

func (s *Service) issueConfirmationToken(ctx context.Context, userID int64) (string, error) {
	token, err := utils.GenerateSecureToken(32)
	if err != nil {
		return "", err
	}
	if err := s.repo.SetConfirmationToken(ctx, userID, token); err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) notifyUserWithToken(email, token string, isResend bool) error {
	return SendConfirmationEmail(email, token, isResend)
}

// HashPassword hashes a password securely using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies a plaintext password against a bcrypt hash.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// ParseFlexibleDuration parses durations like "2d", "3w", or "1mo" into time.Duration.
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
		return time.ParseDuration(input)
	}
}
