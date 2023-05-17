package artists

import artistModels "main/models/artists"

type ReleaseStatsServer interface {
	Get(artistIds []int) artistModels.ArtistReleaseStats
	GetOne(artistId int) artistModels.ArtistReleaseStats
}
