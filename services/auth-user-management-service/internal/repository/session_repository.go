package repository

import (
	"errors"
	"time"

	"auth-user-management-service/internal/domain"

	"gorm.io/gorm"
)

// SessionRepository defines persistence operations for session state.
type SessionRepository interface {
	CreateSession(session *domain.Session) error
	FindByTokenID(tokenID string) (*domain.Session, error)
	RevokeByTokenID(tokenID string) error
}

type sessionRepositoryImpl struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository instance.
func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepositoryImpl{db: db}
}

func (r *sessionRepositoryImpl) CreateSession(session *domain.Session) error {
	return r.db.Create(session).Error
}

func (r *sessionRepositoryImpl) FindByTokenID(tokenID string) (*domain.Session, error) {
	var session domain.Session
	err := r.db.Where("token_id = ?", tokenID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *sessionRepositoryImpl) RevokeByTokenID(tokenID string) error {
	now := time.Now()
	result := r.db.Model(&domain.Session{}).
		Where("token_id = ? AND revoked_at IS NULL", tokenID).
		Update("revoked_at", now)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("session not found or already revoked")
	}
	return nil
}
