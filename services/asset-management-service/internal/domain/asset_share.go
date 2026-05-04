package domain

import "time"

// AssetShare stores explicit read/write sharing grants on folders or notes.
type AssetShare struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"shareId"`
	AssetType        string    `gorm:"type:varchar(20);not null;index:idx_asset_share_lookup,priority:1;index:idx_asset_share_unique,priority:1,unique" json:"assetType"` // folder or note
	AssetID          uint      `gorm:"not null;index:idx_asset_share_lookup,priority:2;index:idx_asset_share_unique,priority:2,unique" json:"assetId"`
	SharedByUserID   uint      `gorm:"not null;index" json:"sharedByUserId"`
	SharedWithUserID uint      `gorm:"not null;index:idx_asset_share_unique,priority:3,unique;index" json:"sharedWithUserId"`
	AccessLevel      string    `gorm:"type:varchar(20);not null" json:"accessLevel"` // read or write
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (AssetShare) TableName() string {
	return "asset_shares"
}
