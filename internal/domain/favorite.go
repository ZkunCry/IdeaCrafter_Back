package domain

import "gorm.io/gorm"

type Favorite struct {
    gorm.Model
    UserID    uint    `gorm:"not null" json:"user_id"`
    StartupID uint    `gorm:"not null" json:"startup_id"`
    User      User    `gorm:"foreignKey:UserID" json:"user"`
    Startup   Startup `gorm:"foreignKey:StartupID" json:"startup"`
}