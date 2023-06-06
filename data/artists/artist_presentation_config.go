package artists

import "time"

// TODO: add artist references
type ArtistPresentationConfig struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null"`
	Value   string    `gorm:"type:jsonb;default:'{}';not null"`
	Updated time.Time `gorm:"not null"`
}
