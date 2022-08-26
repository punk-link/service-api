package artists

import "time"

type Artist struct {
	Id           int       `gorm:"primaryKey,autoIncrement"`
	Created      time.Time `gorm:"not null"`
	ImageDetails string    `gorm:"type:jsonb;default:'{}';not null"`
	LabelId      int       `gorm:"uniqueIndex,not null"`
	Name         string    `gorm:"not null"`
	SpotifyId    string    `gorm:"uniqueIndex,not null"`
	Updated      time.Time `gorm:"not null"`
}
