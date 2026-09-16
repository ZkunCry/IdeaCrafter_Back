package entity

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID    uint    `gorm:"not null;uniqueIndex:idx_favorite_user_startup" json:"user_id"`
	StartupID uint    `gorm:"not null;uniqueIndex:idx_favorite_user_startup" json:"startup_id"`
	User      User    `gorm:"foreignKey:UserID" json:"user"`
	Startup   Startup `gorm:"foreignKey:StartupID" json:"startup"`
}
