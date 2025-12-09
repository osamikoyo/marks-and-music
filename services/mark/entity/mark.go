package entity

import (
	"time"

	"github.com/osamikoyo/music-and-marks/services/mark/api/proto/gen/pb"
)

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

func (m *Mark) ToPB() *pb.Mark {
	return &pb.Mark{
		Id:        uint64(m.ID),
		ReleaseId: m.ReleaseID,
		Value:     m.Value,
		Reviews:   int32(m.Reviews),
	}
}
