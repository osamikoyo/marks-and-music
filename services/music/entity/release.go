package entity

import "time"

type Release struct {
    ID          string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
    MBID        string    `gorm:"uniqueIndex;size:36;not null" json:"mbid"`
    Title       string    `gorm:"type:text;not null" json:"title"`
    ReleaseGroupID string `gorm:"type:uuid;index" json:"release_group_id"`
    ReleaseGroup   ReleaseGroup `gorm:"foreignKey:ReleaseGroupID;references:ID" json:"release_group,omitempty"`
    Status      string    `gorm:"type:text" json:"status,omitempty"`
    Country     string    `gorm:"size:2" json:"country,omitempty"`
    Date        *string   `gorm:"type:date" json:"date,omitempty"`
    Format      string    `gorm:"type:text" json:"format,omitempty"` 
    TrackCount  int       `gorm:"default:0" json:"track_count,omitempty"`
    CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
}