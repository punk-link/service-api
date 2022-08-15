package artists

import "time"

type Artist struct {
	Id        int       `gorm:"primaryKey,autoIncrement"`
	Created   time.Time `gorm:"not null"`
	Name      string    `gorm:"not null"`
	SpotifyId string    `gorm:"not null,index"`
	Updated   time.Time `gorm:"not null"`

	LabelId  int        `gorm:"index"`
	Releases []*Release `gorm:"many2many:release_artists;"`
}
