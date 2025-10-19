package domain

import "gorm.io/gorm"

type Startup struct {
    gorm.Model
    Name        string             `gorm:"not null" json:"name"`
    Description string             `gorm:"type:text;not null" json:"description"`
    ShortDescription string         `gorm:"type:varchar(255);not null" json:"short_description"`
    TargetAudience   string        `gorm:"type:text" json:"target_audience"`
    Problem          string        `gorm:"type:text" json:"problem"`
	Solution         string        `gorm:"type:text" json:"solution"`

    StageID          uint           `json:"stage_id"`                    
	Stage            Stage          `gorm:"foreignKey:StageID" json:"stage"`
    CreatorID   uint                `gorm:"not null" json:"creator_id"`
    Creator     User                `gorm:"foreignKey:CreatorID" json:"creator"`
    Categories  []Category          `gorm:"many2many:startup_categories" json:"categories"` 
    Files       []StartupFile       `gorm:"foreignKey:StartupID;constraint:OnDelete:CASCADE;" json:"files"`
    Vacancies []Vacancy             `gorm:"foreignKey:StartupID;constraint:OnDelete:CASCADE;" json:"vacancies"`

}

