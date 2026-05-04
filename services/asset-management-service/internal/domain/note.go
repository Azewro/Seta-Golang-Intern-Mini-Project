package domain

import (
	"time"

	"gorm.io/gorm"
)

// Note stores text assets within folders.
type Note struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"noteId"`
	FolderID    uint           `gorm:"not null;index" json:"folderId"`
	OwnerUserID uint           `gorm:"not null;index" json:"ownerUserId"`
	Title       string         `gorm:"type:varchar(150);not null" json:"title"`
	Content     string         `gorm:"type:text;not null;default:''" json:"content"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Note) TableName() string {
	return "notes"
}
