package repository

import (
	"auth-user-management-service/internal/domain"

	"gorm.io/gorm"
)

// UserRepository defines the database operations for User
type UserRepository interface {
	CreateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	// Find the first record matching the email. Returns gorm.ErrRecordNotFound if none found.
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
