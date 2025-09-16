package domain

import (
	"gorm.io/gorm"
)

type Application struct {
    gorm.Model
    VacancyID uint `gorm:"not null" json:"vacancy_id"`
    Vacancy   Vacancy `gorm:"foreignKey:VacancyID" json:"vacancy"` 
    UserID    uint `gorm:"not null" json:"user_id"` 
    User      User   `gorm:"foreignKey:UserID" json:"user"` 
    Message   string `json:"message"` 
    Status    string `gorm:"default:'pending'" json:"status"` 
}