package entity

import (
	"time"

	"github.com/osamikoyo/music-and-marks/services/music/api/proto/gen/pb"
)

type Artist struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"type:text;not null" json:"name"`
	SortName  string    `gorm:"type:text" json:"sort_name,omitempty"`
	Country   string    `gorm:"size:2" json:"country,omitempty"`
	Type      string    `gorm:"type:text" json:"type,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

func (a *Artist) ToPB() *pb.Artist {
	return &pb.Artist{
		Id:       a.ID,
		Country:  &a.Country,
		Name:     a.Name,
		SortName: &a.SortName,
		Type:     &a.Type,
	}
}
