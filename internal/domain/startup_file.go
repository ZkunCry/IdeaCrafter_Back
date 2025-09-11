package domain

import "gorm.io/gorm"

type StartupFile struct {
    gorm.Model
    StartupID uint   `gorm:"not null" json:"startup_id"`
    FilePath  string `gorm:"not null" json:"file_path"` 
    FileName  string `gorm:"not null" json:"file_name"`
    MimeType  string `gorm:"type:varchar(100)" json:"mime_type"`
}