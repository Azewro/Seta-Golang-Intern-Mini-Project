package repository

import (
	"time"

	"auth-user-management-service/internal/domain"

	"gorm.io/gorm"
)

// EmailVerificationRepository defines persistence operations for email verification tokens.
type EmailVerificationRepository interface {
	CreateToken(token *domain.EmailVerificationToken) error
	FindByTokenHash(tokenHash string) (*domain.EmailVerificationToken, error)
	MarkUsed(id uint) error
	DeleteActiveTokensByUserID(userID uint) error
	PurgeExpired(now time.Time) error
}

type emailVerificationRepositoryImpl struct {
	db *gorm.DB
}

// NewEmailVerificationRepository creates a new verification repository instance.
func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepositoryImpl{db: db}
}

func (r *emailVerificationRepositoryImpl) CreateToken(token *domain.EmailVerificationToken) error {
	return r.db.Create(token).Error
}

func (r *emailVerificationRepositoryImpl) FindByTokenHash(tokenHash string) (*domain.EmailVerificationToken, error) {
	var token domain.EmailVerificationToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *emailVerificationRepositoryImpl) MarkUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&domain.EmailVerificationToken{}).
		Where("id = ? AND used_at IS NULL", id).
		Update("used_at", now).Error
}

func (r *emailVerificationRepositoryImpl) DeleteActiveTokensByUserID(userID uint) error {
	return r.db.Where("user_id = ? AND used_at IS NULL", userID).
		Delete(&domain.EmailVerificationToken{}).Error
}

func (r *emailVerificationRepositoryImpl) PurgeExpired(now time.Time) error {
	return r.db.Where("expires_at < ?", now).
		Delete(&domain.EmailVerificationToken{}).Error
}
