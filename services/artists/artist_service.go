package artists

import (
	"main/models/artists"
	"main/models/labels"
	"main/services/spotify"
)

type ArtistService struct {
	spotifyService *spotify.SpotifyService
}

func BuildArtistService(spotifyService *spotify.SpotifyService) *ArtistService {
	return &ArtistService{
		spotifyService: spotifyService,
	}
}

func (t *ArtistService) GetRelease(currentManager labels.ManagerContext, id int) artists.Release {
	return artists.Release{}
}

func (t *ArtistService) GetReleases(currentManager labels.ManagerContext, artistId int) []artists.Release {
	return make([]artists.Release, 0)
}

func (t *ArtistService) SearchArtist(query string) []artists.ArtistSearchResult {
	var result []artists.ArtistSearchResult

	const minimalQueryLength int = 3
	if len(query) < minimalQueryLength {
		return result
	}

	results := t.spotifyService.SearchArtist(query)

	return spotify.ToArtistSearchResults(results)
}
