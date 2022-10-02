package platforms

import "time"

type PlatformReleaseUrl struct {
	Id           int       `gorm:"primaryKey,autoIncrement"`
	Created      time.Time `gorm:"not null"`
	PlatformName string    `gorm:"not null"`
	ReleaseId    int       `gorm:"index,not null"`
	Updated      time.Time `gorm:"not null"`
	Url          string    `gorm:"not null"`
}
