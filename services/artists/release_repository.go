package artists

import (
	artistData "main/data/artists"
	"time"
)

type ReleaseRepository interface {
	CreateInBatches(err error, releases *[]artistData.Release) error
	Get(err error, artistId int) ([]artistData.Release, error)
	GetCount(err error) (int64, error)
	GetOne(err error, id int) (artistData.Release, error)
	GetSlimByArtistId(err error, artistId int) ([]artistData.SlimRelease, error)
	GetUpcContainers(err error, top int, skip int, updateTreshold time.Time) ([]artistData.Release, error)
	MarksAsUpdated(err error, ids []int, timestamp time.Time) error
}
