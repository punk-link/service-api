package organizations

import "time"

type Organization struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null"`
	Name    string    `gorm:"not null"`
	Updated time.Time `gorm:"not null"`

	Managers []Manager
}
