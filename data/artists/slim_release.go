package artists

import (
	"time"
)

type SlimRelease struct {
	Id           int       `gorm:"primaryKey,autoIncrement"`
	ImageDetails string    `gorm:"type:jsonb;default:'{}';not null"`
	Name         string    `gorm:"not null"`
	ReleaseDate  time.Time `gorm:"not null"`
	Type         string    `gorm:"not null"`
}
