package spotify

import (
	"fmt"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/responses"
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

func (service *SpotifyService) GetArtistRelease(spotifyId string) responses.ArtistRelease {
	var result responses.ArtistRelease

	var spotifyRelease releases.ArtistRelease
	if err := MakeRequest(service.logger, "GET", fmt.Sprintf("albums/%s", spotifyId), &spotifyRelease); err != nil {
		fmt.Println(err)
		return result
	}

	return toRelease(spotifyRelease)
}

func (service *SpotifyService) GetArtistReleases(spotifyId string) []responses.ArtistRelease {
	spotifyReleases := service.getReleases(spotifyId)
	return toReleases(spotifyReleases)
}

func (service *SpotifyService) SearchArtist(query string) []responses.ArtistSearchResult {
	var result []responses.ArtistSearchResult

	const minimalQueryLength int = 3
	if len(query) < minimalQueryLength {
		return result
	}

	var spotifyArtistSearchResults search.ArtistSearchResult
	err := MakeRequest(service.logger, "GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)), &spotifyArtistSearchResults)
	if err != nil {
		service.logger.LogWarn(err.Error())
		return result
	}

	return toArtistSearchResults(spotifyArtistSearchResults.Artists.Items)
}

func (service *SpotifyService) getReleases(spotifyId string) []releases.ArtistRelease {
	const queryLimit int = 20

	var spotifyResponse releases.ArtistReleaseResult
	offset := 0
	for {
		urlPattern := fmt.Sprintf("artists/%s/albums?limit=%d&offset=%d", spotifyId, queryLimit, offset)
		offset = offset + queryLimit

		var tmpResponse releases.ArtistReleaseResult
		if err := MakeRequest(service.logger, "GET", urlPattern, &tmpResponse); err != nil {
			service.logger.LogWarn(err.Error())
			continue
		}

		spotifyResponse.Items = append(spotifyResponse.Items, tmpResponse.Items...)
		if tmpResponse.Next == "" {
			break
		}
	}

	return spotifyResponse.Items
}
