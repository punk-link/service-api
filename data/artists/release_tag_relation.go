package artists

type ReleaseTagRelation struct {
	ReleaseId int `gorm:"not null,index"`
	TagId     int `gorm:"not null,index"`
}
