package labels

import (
	"main/data/artists"
	"time"
)

type Label struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null"`
	Name    string    `gorm:"not null"`
	Updated time.Time `gorm:"not null"`

	Artists  []artists.Artist
	Managers []Manager
}
