package organizations

import "time"

type Manager struct {
	Id      int       `gorm:"primaryKey,autoIncrement"`
	Created time.Time `gorm:"not null,autoCreateTime" sql:"DEFAULT:now() at time zone 'utc'"`
	Name    string    `gorm:"not null"`
	Updated time.Time `gorm:"not null,autoCreateTime,autoUpdateTime" sql:"DEFAULT:now() at time zone 'utc'"`
}
