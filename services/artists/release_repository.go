package artists

import (
	artistData "main/data/artists"
	"time"
)

type ReleaseRepository interface {
	AddTags(err error, relations *[]artistData.ReleaseTagRelation) error
	CreateInBatches(err error, releases *[]artistData.Release) error
	Get(err error, artistId int) ([]artistData.Release, error)
	GetCount(err error) (int64, error)
	GetOne(err error, id int) (artistData.Release, error)
	GetSlimByArtistId(err error, artistId int) ([]artistData.SlimRelease, error)
	GetTags(err error, releaseIds []int) (map[int][]artistData.Tag, error)
	GetUpcContainers(err error, top int, skip int, updateTreshold time.Time) ([]artistData.Release, error)
	MarksAsUpdated(err error, ids []int, timestamp time.Time) error
}
