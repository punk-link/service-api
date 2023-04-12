package artists

type Tag struct {
	Id               int    `gorm:"primaryKey,autoIncrement"`
	Name             string `gorm:"not null"`
	NormalizedName   string `gorm:"not null"`
	NormalizedVector string `gorm:"not null;type:tsvector"`
}
