package domain

import (
	"time"

	"gorm.io/gorm"
)

// Team stores team ownership and metadata.
type Team struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"teamId"`
	TeamName          string         `gorm:"type:varchar(120);not null" json:"teamName"`
	MainManagerUserID uint           `gorm:"not null;index" json:"mainManagerUserId"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Team) TableName() string {
	return "teams"
}

