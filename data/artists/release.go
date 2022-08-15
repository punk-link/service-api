package artists

import "time"

type Release struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null"`
	//ImageMetadata []spotify.ImageMetadata `json:"images"`
	Label           string    `json:"label"`
	Name            string    `gorm:"not null"`
	PrimaryArtistId int       `gorm:"not null,index"`
	ReleaseDate     time.Time `gorm:"not null"`
	SpotifyId       string    `gorm:"not null,index"`
	TrackNumber     int       `gorm:"not null"`
	//Tracks          TrackContainer `json:"tracks"`
	Type    string    `gorm:"not null"`
	Updated time.Time `gorm:"not null"`

	Artists []*Artist `gorm:"many2many:release_artists;"`
}
