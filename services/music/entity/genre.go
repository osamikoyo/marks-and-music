package entity

import "time"

type Genre struct {
    ID    uint   `gorm:"primaryKey" json:"id"`
    Name  string `gorm:"size:100;not null;uniqueIndex" json:"name"`
    
    Albums []Album `gorm:"foreignKey:GenreID" json:"-"`
    
    CreatedAt time.Time `json:"created_at,omitempty"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`
}