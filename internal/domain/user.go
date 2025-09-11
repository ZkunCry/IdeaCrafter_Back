package domain

import "gorm.io/gorm"

type User struct {
    gorm.Model
    Username     string    `gorm:"unique;not null" json:"username"`
    Email        string    `gorm:"unique;not null" json:"email"`
    PasswordHash string    `gorm:"not null" json:"-"` 
    Startups     []Startup `gorm:"foreignKey:CreatorID" json:"startups"` 
    Favorites    []Favorite `gorm:"foreignKey:UserID" json:"favorites"`
    Memberships  []Membership `gorm:"foreignKey:UserID" json:"memberships"`
}