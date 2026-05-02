package repository

import (
	"team-management-service/internal/domain"

	"gorm.io/gorm"
)

// UserRepository provides read-only access to global user records.
type UserRepository interface {
	FindByID(id uint) (*domain.UserRef, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindByID(id uint) (*domain.UserRef, error) {
	var user domain.UserRef
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
