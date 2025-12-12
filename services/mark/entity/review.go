package entity

import (
	"time"

	"github.com/osamikoyo/music-and-marks/services/mark/api/proto/gen/pb"
)

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

func (r *Review) ToPB() *pb.Review {
	return &pb.Review{
		Id:        uint64(r.ID),
		Text:      r.Text,
		Count:     int32(r.Count),
		UserId:    r.UserID,
		ReleaseId: r.ReleaseID,
	}
}
