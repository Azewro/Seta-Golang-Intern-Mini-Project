package repository

import (
	"auth-user-management-service/internal/domain"

	"gorm.io/gorm"
)

// UserRepository defines the database operations for User
type UserRepository interface {
	CreateUser(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	ListUsers(offset int, limit int) ([]domain.User, error)
	MarkUserVerified(id uint) error
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
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) ListUsers(offset int, limit int) ([]domain.User, error) {
	var users []domain.User
	err := r.db.Offset(offset).Limit(limit).Order("id ASC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) MarkUserVerified(id uint) error {
	return r.db.Model(&domain.User{}).Where("id = ?", id).Updates(map[string]any{
		"is_verified": true,
	}).Error
}
