package entity

import (
	"time"
)

type Artist struct {
    ID             string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    Name           string    `gorm:"type:text;not null" json:"name"`
    SortName       string    `gorm:"type:text" json:"sort_name,omitempty"`
    Country        string    `gorm:"size:2" json:"country,omitempty"`
    Type           string    `gorm:"type:text" json:"type,omitempty"` 
    CreatedAt      time.Time `gorm:"autoCreateTime" json:"-"`
    UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"-"`
}