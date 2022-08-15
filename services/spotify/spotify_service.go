package spotify

import (
	"fmt"
	"main/models/artists"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/services/common"
	"net/url"
	"strings"
)

type SpotifyService struct {
	logger *common.Logger
}

func BuildSpotifyService(logger *common.Logger) *SpotifyService {
	return &SpotifyService{
		logger: logger,
	}
}

func (t *SpotifyService) GetReleaseDetails(spotifyId string) artists.Release {
	var result artists.Release

	var spotifyRelease releases.Release
	if err := makeRequest(t.logger, "GET", fmt.Sprintf("albums/%s", spotifyId), &spotifyRelease); err != nil {
		t.logger.LogWarn(err.Error())
		return result
	}

	return toRelease(spotifyRelease)
}

func (t *SpotifyService) GetReleasesDetails(spotifyIds []string) []releases.Release {
	const queryLimit int = 20

	spotifyReleases := make([]releases.Release, 0)
	skip := 0
	for i := 0; i < len(spotifyIds); i = i + queryLimit {
		ids := spotifyIds[skip:getSliceEnd(&spotifyIds, skip+queryLimit)]
		idQuery := strings.Join(ids, ",")

		var tmpResponse releases.ReleaseDetailsContainer
		if err := makeRequest(t.logger, "GET", fmt.Sprintf("albums?ids=%s", idQuery), &tmpResponse); err != nil {
			t.logger.LogWarn(err.Error())
			return spotifyReleases
		}

		spotifyReleases = append(spotifyReleases, tmpResponse.Items...)
		skip += queryLimit
	}

	return spotifyReleases
}

func (t *SpotifyService) GetArtistReleases(spotifyId string) []releases.Release {
	const queryLimit int = 20

	var spotifyResponse releases.ReleaseContainer
	offset := 0
	for {
		urlPattern := fmt.Sprintf("artists/%s/albums?limit=%d&offset=%d", spotifyId, queryLimit, offset)
		offset = offset + queryLimit

		var tmpResponse releases.ReleaseContainer
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

func (t *SpotifyService) SearchArtist(query string) []search.Artist {
	var spotifyArtistSearchResults search.ArtistSearchResult
	err := makeRequest(t.logger, "GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)), &spotifyArtistSearchResults)
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]search.Artist, 0)
	}

	return spotifyArtistSearchResults.Artists.Items
}

func getSliceEnd[T any](slice *[]T, iterationEnd int) int {
	sliceEnd := len(*slice)

	if sliceEnd < iterationEnd {
		return sliceEnd
	}

	return iterationEnd
}
