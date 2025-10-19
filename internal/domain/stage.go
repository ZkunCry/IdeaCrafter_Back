package domain

import "gorm.io/gorm"

type Stage struct {
	gorm.Model
	Name     string     `gorm:"unique;not null" json:"name"`
	Startups []Startup  `gorm:"foreignKey:StageID" json:"startups"`
}

