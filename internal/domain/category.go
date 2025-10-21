package domain

import "gorm.io/gorm"

type Category struct {
    gorm.Model
    Name     string     `gorm:"unique;not null" json:"name"` 
    Startups []Startup  `gorm:"many2many:startup_categories" json:"-"`
}