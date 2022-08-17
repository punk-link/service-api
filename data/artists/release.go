package artists

import (
	"time"

	"github.com/jackc/pgtype"
)

type Release struct {
	Id              int               `gorm:"primaryKey,autoIncrement"`
	Created         time.Time         `gorm:"not null"`
	ImageDetails    pgtype.JSONB      `gorm:"type:jsonb;default:'{}';not null"`
	Label           string            `gorm:"not null"`
	Name            string            `gorm:"not null"`
	PrimaryArtistId int               `gorm:"not null,index"`
	ReleaseDate     time.Time         `gorm:"not null"`
	SpotifyId       string            `gorm:"not null,index"`
	TrackNumber     int               `gorm:"not null"`
	Tracks          pgtype.JSONBArray `gorm:"type:jsonb;default:'[]';not null"`
	Type            string            `gorm:"not null"`
	Updated         time.Time         `gorm:"not null"`

	Artists []*Artist `gorm:"many2many:release_artists;"`
}
