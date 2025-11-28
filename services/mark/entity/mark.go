package entity

import "time"

type Mark struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ReleaseID string    `json:"release_id"`
	Value     float32   `json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
