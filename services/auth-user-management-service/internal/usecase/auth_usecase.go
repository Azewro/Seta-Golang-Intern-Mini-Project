package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"auth-user-management-service/internal/domain"
	"auth-user-management-service/internal/repository"
	"auth-user-management-service/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyUsed   = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrSessionUnavailable = errors.New("session not found or already revoked")
	ErrJWTSecretMissing   = errors.New("jwt secret is not configured")
	ErrEmailNotVerified   = errors.New("email is not verified")
	ErrInvalidVerifyToken = errors.New("invalid verification token")
	ErrExpiredVerifyToken = errors.New("verification token expired")
	ErrSMTPNotConfigured  = errors.New("smtp is not fully configured")
)

// RegisterRequest defines the input payload for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"omitempty,oneof=manager member"`
}

// LoginRequest defines the input payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// ResendVerificationRequest defines payload to resend verification emails.
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// BulkGetUsersRequest resolves a batch of users by ID (used by other services).
type BulkGetUsersRequest struct {
	UserIDs []uint `json:"userIds" binding:"required"`
}

// UserResponse is a safe user payload without password.
type UserResponse struct {
	ID         uint      `json:"userId"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// LoginResponse defines login output payload.
type LoginResponse struct {
	AccessToken string       `json:"accessToken"`
	ExpiresAt   time.Time    `json:"expiresAt"`
	User        UserResponse `json:"user"`
}

// AuthUsecase defines the business logic interface.
type AuthUsecase interface {
	Register(req *RegisterRequest) (*domain.User, error)
	VerifyEmail(rawToken string) error
	ResendVerification(req *ResendVerificationRequest) error
	Login(req *LoginRequest) (*LoginResponse, error)
	Logout(tokenID string) error
	ListUsers(page int, limit int) ([]UserResponse, error)
	GetMyProfile(userID uint) (*UserResponse, error)
	GetUsersByIDs(ids []uint) ([]UserResponse, error)
}

type authUsecaseImpl struct {
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
	verifyRepo     repository.EmailVerificationRepository
	jwtSecret      string
	jwtTokenTTL    time.Duration
	verifyTokenTTL time.Duration
	appBaseURL     string
	smtpConfig     utils.SMTPConfig
}

// NewAuthUsecase creates a new usecase instance (Dependency Injection).
func NewAuthUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	verifyRepo repository.EmailVerificationRepository,
) AuthUsecase {
	ttl := 24 * time.Hour
	jwtExpiresHours := strings.TrimSpace(os.Getenv("JWT_EXPIRES_HOURS"))
	if jwtExpiresHours != "" {
		if parsed, err := time.ParseDuration(jwtExpiresHours + "h"); err == nil {
			ttl = parsed
		}
	}

	verifyTTLMinutes := readIntEnvWithDefault("EMAIL_VERIFY_TOKEN_TTL_MINUTES", 5)

	appBaseURL := strings.TrimSuffix(strings.TrimSpace(os.Getenv("APP_BASE_URL")), "/")
	if appBaseURL == "" {
		appBaseURL = "http://localhost:8080"
	}

	return &authUsecaseImpl{
		userRepo:       userRepo,
		sessionRepo:    sessionRepo,
		verifyRepo:     verifyRepo,
		jwtSecret:      strings.TrimSpace(os.Getenv("JWT_SECRET")),
		jwtTokenTTL:    ttl,
		verifyTokenTTL: time.Duration(verifyTTLMinutes) * time.Minute,
		appBaseURL:     appBaseURL,
		smtpConfig:     utils.LoadSMTPConfig(),
	}
}

// Register processes user registration and sends a verification link.
func (u *authUsecaseImpl) Register(req *RegisterRequest) (*domain.User, error) {
	existingUser, err := u.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "member"
	}

	user := &domain.User{
		Username:   req.Username,
		Email:      req.Email,
		Password:   hashedPassword,
		Role:       role,
		IsVerified: false,
	}

	if err := u.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	if err := u.issueVerification(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *authUsecaseImpl) VerifyEmail(rawToken string) error {
	tokenHash := hashToken(rawToken)
	verificationToken, err := u.verifyRepo.FindByTokenHash(tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidVerifyToken
		}
		return err
	}

	if verificationToken.UsedAt != nil {
		return ErrInvalidVerifyToken
	}

	if time.Now().After(verificationToken.ExpiresAt) {
		return ErrExpiredVerifyToken
	}

	if err := u.userRepo.MarkUserVerified(verificationToken.UserID); err != nil {
		return err
	}
	if err := u.verifyRepo.MarkUsed(verificationToken.ID); err != nil {
		return err
	}

	return nil
}

func (u *authUsecaseImpl) ResendVerification(req *ResendVerificationRequest) error {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	if user.IsVerified {
		return nil
	}

	if err := u.verifyRepo.PurgeExpired(time.Now()); err != nil {
		return err
	}

	return u.issueVerification(user)
}

func (u *authUsecaseImpl) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	if !user.IsVerified {
		return nil, ErrEmailNotVerified
	}

	if u.jwtSecret == "" {
		return nil, ErrJWTSecretMissing
	}

	tokenID := uuid.NewString()
	token, expiresAt, err := utils.GenerateToken(user.ID, user.Role, tokenID, u.jwtSecret, u.jwtTokenTTL)
	if err != nil {
		return nil, err
	}

	err = u.sessionRepo.CreateSession(&domain.Session{
		UserID:    user.ID,
		TokenID:   tokenID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
		User:        toUserResponse(user),
	}, nil
}

func (u *authUsecaseImpl) Logout(tokenID string) error {
	err := u.sessionRepo.RevokeByTokenID(tokenID)
	if err != nil {
		return ErrSessionUnavailable
	}
	return nil
}

func (u *authUsecaseImpl) ListUsers(page int, limit int) ([]UserResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
	users, err := u.userRepo.ListUsers(offset, limit)
	if err != nil {
		return nil, err
	}

	response := make([]UserResponse, 0, len(users))
	for i := range users {
		response = append(response, toUserResponse(&users[i]))
	}

	return response, nil
}

func (u *authUsecaseImpl) GetMyProfile(userID uint) (*UserResponse, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	response := toUserResponse(user)
	return &response, nil
}

func (u *authUsecaseImpl) GetUsersByIDs(ids []uint) ([]UserResponse, error) {
	users, err := u.userRepo.FindUsersByIDs(ids)
	if err != nil {
		return nil, err
	}

	response := make([]UserResponse, 0, len(users))
	for i := range users {
		response = append(response, toUserResponse(&users[i]))
	}

	return response, nil
}

func (u *authUsecaseImpl) issueVerification(user *domain.User) error {
	if !isSMTPConfigured(u.smtpConfig) {
		return ErrSMTPNotConfigured
	}

	if err := u.verifyRepo.PurgeExpired(time.Now()); err != nil {
		return err
	}

	if err := u.verifyRepo.DeleteActiveTokensByUserID(user.ID); err != nil {
		return err
	}

	rawToken := uuid.NewString()
	tokenHash := hashToken(rawToken)
	expiresAt := time.Now().Add(u.verifyTokenTTL)

	verificationToken := &domain.EmailVerificationToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}

	if err := u.verifyRepo.CreateToken(verificationToken); err != nil {
		return err
	}

	verifyURL := fmt.Sprintf("%s/api/v1/auth/verify-email?token=%s", u.appBaseURL, rawToken)
	if err := utils.SendVerificationEmail(u.smtpConfig, user.Email, verifyURL); err != nil {
		return err
	}

	return nil
}

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Role:       user.Role,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func hashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(rawToken)))
	return hex.EncodeToString(sum[:])
}

func readIntEnvWithDefault(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func isSMTPConfigured(cfg utils.SMTPConfig) bool {
	return cfg.Host != "" && cfg.Username != "" && cfg.Password != ""
}
