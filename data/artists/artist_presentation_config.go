package artists

import "time"

type ArtistPresentationConfig struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null"`
	Value   string    `gorm:"type:jsonb;default:'{}';not null"`
	Updated time.Time `gorm:"not null"`
}
