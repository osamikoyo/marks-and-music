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

func NewReview(releaeID, text, userID string, count int) *Review {
	return &Review{
		Text:      text,
		ReleaseID: releaeID,
		UserID:    userID,
		Count:     count,
		CreatedAt: time.Now(),
	}
}
