package artists

import (
	"time"
)

type Release struct {
	Id                 int       `gorm:"primaryKey,autoIncrement"`
	Created            time.Time `gorm:"not null"`
	FeaturingArtistIds string    `gorm:"type:jsonb;default:'[]';not null"`
	ImageDetails       string    `gorm:"type:jsonb;default:'{}';not null"`
	Label              string    `gorm:"not null"`
	Name               string    `gorm:"not null"`
	ReleaseArtistIds   string    `gorm:"type:jsonb;default:'[]';not null"`
	ReleaseDate        time.Time `gorm:"not null"`
	SpotifyId          string    `gorm:"not null,uniqueIndex"`
	TrackNumber        int       `gorm:"not null"`
	Tracks             string    `gorm:"type:jsonb;default:'[]';not null"`
	Type               string    `gorm:"not null"`
	Upc                string    `gorm:"not null"`
	Updated            time.Time `gorm:"not null"`
}
