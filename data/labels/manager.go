package labels

import "time"

type Manager struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null"`
	Name    string    `gorm:"not null"`
	Updated time.Time `gorm:"not null"`

	LabelId int
}
