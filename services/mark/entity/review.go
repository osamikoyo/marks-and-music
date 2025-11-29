package entity

import "time"

type Review struct {
	ID        uint      `gorm:"primarKey" json:"id"`
	Text      string    `json:"text"`
	Count     int       `json:"count"`
	UserID    string    `json:"user_id"`
	ReleaseID string    `json:"release_id"`
	CreatedAt time.Time `json:"-"`
}
