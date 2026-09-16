package entity

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name         string    `gorm:"unique;not null" json:"name"`
	Slug         string    `gorm:"unique;default:''" json:"slug"`
	StartupCount int       `gorm:"column:startup_count;->" json:"-"`
	Startups     []Startup `gorm:"many2many:startup_categories" json:"-"`
}
