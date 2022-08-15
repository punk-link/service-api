package artists

import (
	"main/models/artists"
	"main/models/labels"
	"main/services/spotify"
)

type ReleaseService struct {
	spotifyService *spotify.SpotifyService
}

func BuildReleaseService(spotifyService *spotify.SpotifyService) *ReleaseService {
	return &ReleaseService{
		spotifyService: spotifyService,
	}
}

func (t *ReleaseService) GetRelease(currentManager labels.ManagerContext, id int) artists.Release {
	return artists.Release{}
}

func (t *ReleaseService) GetReleases(currentManager labels.ManagerContext, artistId int) []artists.Release {
	return make([]artists.Release, 0)
}
