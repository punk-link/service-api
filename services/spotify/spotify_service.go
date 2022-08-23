package spotify

import (
	"fmt"
	"main/helpers"
	spotifyArtists "main/models/spotify/artists"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/services/common"
	"net/url"
	"strings"
)

type SpotifyService struct {
	logger *common.Logger
}

func ConstructSpotifyService(logger *common.Logger) *SpotifyService {
	return &SpotifyService{
		logger: logger,
	}
}

func (t *SpotifyService) GetArtist(spotifyId string) (spotifyArtists.Artist, error) {
	var spotifyArtist spotifyArtists.Artist
	if err := makeRequest(t.logger, "GET", fmt.Sprintf("artists/%s", spotifyId), &spotifyArtist); err != nil {
		t.logger.LogWarn(err.Error())
		return spotifyArtists.Artist{}, err
	}

	return spotifyArtist, nil
}

func (t *SpotifyService) GetArtists(spotifyIds []string) []spotifyArtists.Artist {
	const queryLimit int = 50
	chunkedIds := helpers.Chunk(spotifyIds, queryLimit)

	urls := make([]string, len(chunkedIds))
	for i, chunk := range chunkedIds {
		ids := strings.Join(chunk, ",")
		urls[i] = fmt.Sprintf("artists?ids=%s", ids)
	}

	spotifyArtistContainers := makeBatchRequest[spotifyArtists.ArtistContainer](t.logger, "GET", urls)

	results := make([]spotifyArtists.Artist, 0)
	for _, container := range spotifyArtistContainers {
		results = append(results, container.Artists...)
	}

	return results
}

func (t *SpotifyService) GetArtistReleases(spotifyId string) []releases.Release {
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

func (t *SpotifyService) GetReleasesDetails(spotifyIds []string) []releases.Release {
	chunkedIds := helpers.Chunk(spotifyIds, queryLimit)
	urls := make([]string, len(chunkedIds))
	for i, chunk := range chunkedIds {
		ids := strings.Join(chunk, ",")
		urls[i] = fmt.Sprintf("albums?ids=%s", ids)
	}

	releaseContainers := makeBatchRequest[releases.ReleaseDetailsContainer](t.logger, "GET", urls)

	spotifyReleases := make([]releases.Release, 0)
	for _, container := range releaseContainers {
		spotifyReleases = append(spotifyReleases, container.Releases...)
	}

	return spotifyReleases
}

func (t *SpotifyService) SearchArtist(query string) []spotifyArtists.SlimArtist {
	var spotifyArtistSearchResults search.ArtistSearchResult
	err := makeRequest(t.logger, "GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)), &spotifyArtistSearchResults)
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]spotifyArtists.SlimArtist, 0)
	}

	return spotifyArtistSearchResults.Artists.Items
}

const queryLimit int = 20
