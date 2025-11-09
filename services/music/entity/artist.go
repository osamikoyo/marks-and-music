package entity

import (
	"time"

	"gorm.io/gorm"
)

type Artist struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"size:255;not null;uniqueIndex" json:"name"`
    Country   string         `gorm:"size:100" json:"country,omitempty"`
    FoundedAt *time.Time     `json:"founded_at,omitempty"`
    
    Albums    []Album        `gorm:"foreignKey:ArtistID" json:"-"`
    Songs     []Song         `gorm:"foreignKey:ArtistID" json:"-"`
    
    CreatedAt time.Time      `json:"created_at,omitempty"`
    UpdatedAt time.Time      `json:"updated_at,omitempty"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}