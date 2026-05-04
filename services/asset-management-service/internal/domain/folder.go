package domain

import (
	"time"

	"gorm.io/gorm"
)

// Folder stores user-owned asset containers.
type Folder struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"folderId"`
	OwnerUserID uint           `gorm:"not null;index" json:"ownerUserId"`
	Name        string         `gorm:"type:varchar(150);not null" json:"name"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Folder) TableName() string {
	return "folders"
}
