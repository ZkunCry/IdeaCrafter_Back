package domain

import (
	"gorm.io/gorm"
)
type Vacancy struct{
	gorm.Model
	StartupID uint `gorm:"null;constraint:OnDelete:CASCADE;" json:"startup_id"`
	RoleID uint 	`gorm:"not null" json:"role_id"`
	Role  Role 		`gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;" json:"role"`
	Description string `json:"description"`
	IsOpen bool `gorm:"default:true" json:"is_open"`
	UserID *uint `gorm:"default:null" json:"user_id"` 
  User   *User `gorm:"foreignKey:UserID" json:"user"`
}
