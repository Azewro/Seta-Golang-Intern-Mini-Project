package usecase

import (
	"errors"

	"auth-user-management-service/internal/domain"
	"auth-user-management-service/internal/repository"
	"auth-user-management-service/pkg/utils"
)

// RegisterRequest defines the input payload for user registration
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role"` // Optional: manager or member
}

// AuthUsecase defines the business logic interface
type AuthUsecase interface {
	Register(req *RegisterRequest) (*domain.User, error)
}

type authUsecaseImpl struct {
	userRepo repository.UserRepository
}

// NewAuthUsecase creates a new usecase instance (Dependency Injection)
func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecaseImpl{
		userRepo: userRepo,
	}
}

// Register processes the user registration logic
func (u *authUsecaseImpl) Register(req *RegisterRequest) (*domain.User, error) {
	// 1. Check if email already exists
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	// 2. Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 3. Assign role (default is member, strict definition)
	role := "member"
	if req.Role == "manager" {
		role = "manager"
	}

	// 4. Create the domain.User entity
	newUser := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     role,
	}

	// 5. Save using repository
	if err := u.userRepo.CreateUser(newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	return newUser, nil
}

