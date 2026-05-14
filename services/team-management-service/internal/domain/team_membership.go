package domain

import "time"

// TeamMembership maps users into teams with a team-scoped role.
type TeamMembership struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID         uint      `gorm:"not null;index:idx_team_user,unique" json:"teamId"`
	UserID         uint      `gorm:"not null;index:idx_team_user,unique;index" json:"userId"`
	MembershipRole string    `gorm:"type:varchar(20);not null" json:"membershipRole"` // manager or member
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (TeamMembership) TableName() string {
	return "team_memberships"
}
