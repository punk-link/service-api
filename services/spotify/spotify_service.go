package spotify

import (
	"fmt"
	"main/models/artists"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/services/common"
	"net/url"
)

type SpotifyService struct {
	logger *common.Logger
}

func BuildSpotifyService(logger *common.Logger) *SpotifyService {
	return &SpotifyService{
		logger: logger,
	}
}

func (t *SpotifyService) GetArtistRelease(spotifyId string) artists.Release {
	var result artists.Release

	var spotifyRelease releases.ArtistRelease
	if err := makeRequest(t.logger, "GET", fmt.Sprintf("albums/%s", spotifyId), &spotifyRelease); err != nil {
		fmt.Println(err)
		return result
	}

	return toRelease(spotifyRelease)
}

func (t *SpotifyService) GetArtistReleases(spotifyId string) []artists.Release {
	spotifyReleases := t.getReleases(spotifyId)
	return toReleases(spotifyReleases)
}

func (t *SpotifyService) SearchArtist(query string) []search.Artist {
	var spotifyArtistSearchResults search.ArtistSearchResult
	err := makeRequest(t.logger, "GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)), &spotifyArtistSearchResults)
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]search.Artist, 0)
	}

	return spotifyArtistSearchResults.Artists.Items
}

func (t *SpotifyService) getReleases(spotifyId string) []releases.ArtistRelease {
	const queryLimit int = 20

	var spotifyResponse releases.ArtistReleaseResult
	offset := 0
	for {
		urlPattern := fmt.Sprintf("artists/%s/albums?limit=%d&offset=%d", spotifyId, queryLimit, offset)
		offset = offset + queryLimit

		var tmpResponse releases.ArtistReleaseResult
		if err := makeRequest(t.logger, "GET", urlPattern, &tmpResponse); err != nil {
			t.logger.LogWarn(err.Error())
			continue
		}

		spotifyResponse.Items = append(spotifyResponse.Items, tmpResponse.Items...)
		if tmpResponse.Next == "" {
			break
		}
	}

	return spotifyResponse.Items
}
