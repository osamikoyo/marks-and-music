package entity

import (
	"time"

	"gorm.io/gorm"
)

type Album struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    Title      string    `gorm:"size:255;not null" json:"title"`
    ArtistID   uint      `gorm:"not null" json:"-"`
    GenreID    uint      `gorm:"not null" json:"-"`
    ReleasedAt time.Time `gorm:"not null" json:"released_at"`
    
    Artist     Artist    `gorm:"foreignKey:ArtistID" json:"artist"`
    Genre      Genre     `gorm:"foreignKey:GenreID" json:"genre"`
    Songs      []Song    `gorm:"foreignKey:AlbumID" json:"songs,omitempty"`
    
    CreatedAt time.Time      `json:"created_at,omitempty"`
    UpdatedAt time.Time      `json:"updated_at,omitempty"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}