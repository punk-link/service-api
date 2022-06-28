package managers

import "time"

type Manager struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null,autoCreateTime"`
	Name    string    `gorm:"not null"`
	Updated time.Time `gorm:"not null,autoCreateTime,autoUpdateTime"`
}
