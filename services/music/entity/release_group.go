package entity

import "time"

type ReleaseGroup struct {
    ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    MBID             string    `gorm:"uniqueIndex;size:36;not null" json:"mbid"`
    Title            string    `gorm:"type:text;not null" json:"title"`
    ArtistID         string    `gorm:"type:uuid;index" json:"artist_id"`
    Artist           Artist    `gorm:"foreignKey:ArtistID;references:ID" json:"artist,omitempty"`
    PrimaryType      string    `gorm:"type:text" json:"primary_type"` // Album, Single, EP...
    SecondaryTypes   []string  `gorm:"type:text[]" json:"secondary_types,omitempty"` // Live, Compilation
    FirstReleaseDate *string   `gorm:"type:date" json:"first_release_date,omitempty"` // "2025-03-14"
    CreatedAt        time.Time `gorm:"autoCreateTime" json:"-"`
    UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"-"`
}