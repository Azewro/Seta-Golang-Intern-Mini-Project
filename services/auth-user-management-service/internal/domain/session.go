package domain

import "time"

// Session stores token state to support real logout/revocation.
type Session struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"userId"`
	TokenID   string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"tokenId"`
	ExpiresAt time.Time  `gorm:"not null" json:"expiresAt"`
	RevokedAt *time.Time `json:"revokedAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
