package artists

import "time"

type Artist struct {
	Id        int       `gorm:"primaryKey,autoIncrement"`
	Created   time.Time `gorm:"not null"`
	LabelId   int       `gorm:"uniqueIndex,not null"`
	Name      string    `gorm:"not null"`
	SpotifyId string    `gorm:"uniqueIndex,not null"`
	Updated   time.Time `gorm:"not null"`

	Releases []Release `gorm:"many2many:release_artists;"`
}
