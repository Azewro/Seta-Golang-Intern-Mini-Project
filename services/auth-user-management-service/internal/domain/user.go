package domain

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user record in the database (Entity/Model)
type User struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"userId"`
	Username   string         `gorm:"type:varchar(100);not null" json:"username"`
	Email      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password   string         `gorm:"type:varchar(255);not null" json:"-"` // Hide password when parse/marshal JSON
	Role       string         `gorm:"type:varchar(50);not null;default:'member'" json:"role"` // manager or member
	IsVerified bool           `gorm:"not null;default:false" json:"isVerified"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
