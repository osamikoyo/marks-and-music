package entity

import "time"

type Mark struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	ReleaseID string
	Value     float32
	Reviews   int
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func NewMark(releaeID string, value float32) *Mark {
	return &Mark{
		ReleaseID: releaeID,
		Value:     value,
		CreatedAt: time.Now(),
	}
}
