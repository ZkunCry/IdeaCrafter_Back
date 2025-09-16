package domain

import "gorm.io/gorm"

type Startup struct {
    gorm.Model
    Name        string       `gorm:"not null" json:"name"`
    Description string       `gorm:"type:text;not null" json:"description"`
    CreatorID   uint         `gorm:"not null" json:"creator_id"`
    Creator     User         `gorm:"foreignKey:CreatorID" json:"creator"`
    Categories  []Category   `gorm:"many2many:startup_categories" json:"categories"` 
    Files       []StartupFile `gorm:"foreignKey:StartupID" json:"files"`

    Vacancies []Vacancy `gorm:"foreignKey:StartupID" json:"vacancies"`
}

type CreateStartupInput struct {
    Name        string   `json:"name" validate:"required"`
    Description string   `json:"description" validate:"required"`
    CategoryIDs []uint   `json:"category_ids" validate:"required"`
}