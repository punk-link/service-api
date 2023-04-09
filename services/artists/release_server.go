package artists

import (
	artistData "main/data/artists"
	artistModels "main/models/artists"
	labelModels "main/models/labels"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"
	"time"

	platformContracts "github.com/punk-link/platform-contracts"
)

type ReleaseServer interface {
	Add(currentManager labelModels.ManagerContext, artists map[string]artistData.Artist, releases []releaseSpotifyPlatformModels.Release, timeStamp time.Time) error
	Get(artistId int) ([]artistModels.Release, error)
	GetCount() int
	GetMissing(artistId int, artistSpotifyId string) ([]releaseSpotifyPlatformModels.Release, error)
	GetOne(id int) (artistModels.Release, error)
	GetUpcContainersToUpdate(top int, skip int, updateTreshold time.Time) []platformContracts.UpcContainer
	MarkAsUpdated(ids []int, timestamp time.Time) error
}
