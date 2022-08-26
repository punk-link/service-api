package artists

type ArtistReleaseRelation struct {
	ArtistId  int `gorm:"not null,index"`
	ReleaseId int `gorm:"not null,index"`
}
