package entity

import (
	"time"

	"gorm.io/gorm"
)

type Song struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    Title    string `gorm:"size:255;not null" json:"title"`
    AlbumID  uint   `gorm:"not null" json:"-"`
    ArtistID uint   `gorm:"not null" json:"-"`
    TrackNo  int    `gorm:"not null" json:"track_number"`
    Duration int    `gorm:"comment:'Длительность в секундах'" json:"duration_seconds,omitempty"`
    
    Album  Album  `gorm:"foreignKey:AlbumID" json:"album"`
    Artist Artist `gorm:"foreignKey:ArtistID" json:"artist"`
    
    CreatedAt time.Time      `json:"created_at,omitempty"`
    UpdatedAt time.Time      `json:"updated_at,omitempty"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}