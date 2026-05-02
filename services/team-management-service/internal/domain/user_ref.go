package domain

// UserRef is a read-only mapping to auth service users table.
type UserRef struct {
	ID   uint   `gorm:"column:id;primaryKey"`
	Role string `gorm:"column:role"`
}

func (UserRef) TableName() string {
	return "users"
}

