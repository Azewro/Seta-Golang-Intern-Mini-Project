package usecase

import (
	"errors"
	"os"
	"strings"
	"time"

	"auth-user-management-service/internal/domain"
	"auth-user-management-service/internal/repository"
	"auth-user-management-service/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyUsed    = errors.New("email already in use")
	ErrInvalidRole         = errors.New("role must be either manager or member")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrSessionUnavailable  = errors.New("session not found or already revoked")
	ErrJWTSecretMissing    = errors.New("jwt secret is not configured")
)

// RegisterRequest defines the input payload for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

// LoginRequest defines the input payload for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is a safe user payload without password.
type UserResponse struct {
	ID        uint      `json:"userId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// LoginResponse defines login output payload.
type LoginResponse struct {
	AccessToken string       `json:"accessToken"`
	ExpiresAt   time.Time    `json:"expiresAt"`
	User        UserResponse `json:"user"`
}

// AuthUsecase defines the business logic interface
type AuthUsecase interface {
	Register(req *RegisterRequest) (*domain.User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	Logout(tokenID string) error
	ListUsers(page int, limit int) ([]UserResponse, error)
	GetMyProfile(userID uint) (*UserResponse, error)
}

type authUsecaseImpl struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	jwtSecret    string
	jwtTokenTTL  time.Duration
}

// NewAuthUsecase creates a new usecase instance (Dependency Injection)
func NewAuthUsecase(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) AuthUsecase {
	ttl := 24 * time.Hour
	jwtExpiresHours := strings.TrimSpace(os.Getenv("JWT_EXPIRES_HOURS"))
	if jwtExpiresHours != "" {
		if parsed, err := time.ParseDuration(jwtExpiresHours + "h"); err == nil {
			ttl = parsed
		}
	}

	return &authUsecaseImpl{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   strings.TrimSpace(os.Getenv("JWT_SECRET")),
		jwtTokenTTL: ttl,
	}
}

// Register processes the user registration logic
func (u *authUsecaseImpl) Register(req *RegisterRequest) (*domain.User, error) {
	// 1. Validate strict role at creation.
	if req.Role != "manager" && req.Role != "member" {
		return nil, ErrInvalidRole
	}

	// 2. Check if email already exists.
	existingUser, err := u.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 3. Hash password.
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 4. Create user entity.
	newUser := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	// 5. Persist user.
	if err := u.userRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
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

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

