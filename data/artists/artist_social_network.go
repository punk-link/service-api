package artists

import "time"

type ArtistSocialNetwork struct {
	Id        int       `gorm:"primaryKey,autoIncrement"`
	ArtistId  int       `gorm:"index,not null"`
	Created   time.Time `gorm:"not null"`
	NetworkId string    `gorm:"index,not null"`
	Url       string    `gorm:"not null"`
	Updated   time.Time `gorm:"not null"`
}
