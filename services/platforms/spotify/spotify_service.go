package spotify

import (
	"fmt"
	"main/helpers"
	spotifyArtists "main/models/platforms/spotify/artists"
	"main/models/platforms/spotify/releases"
	searchModels "main/models/platforms/spotify/searches"
	spotifyModels "main/models/platforms/spotify/tokens"
	"net/http"
	"net/url"
	"strings"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type SpotifyService struct {
	config     *spotifyModels.SpotifyClientConfig
	httpConfig *httpClient.HttpClientConfig
	logger     logger.Logger
}

func NewSpotifyService(injector *do.Injector) (*SpotifyService, error) {
	config := do.MustInvoke[*spotifyModels.SpotifyClientConfig](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	httpConfig := httpClient.DefaultConfig(logger)

	return &SpotifyService{
		config:     config,
		httpConfig: httpConfig,
		logger:     logger,
	}, nil
}

func (t *SpotifyService) GetArtist(spotifyId string) (spotifyArtists.Artist, error) {
	httpRequest, err := t.getHttpRequest("GET", fmt.Sprintf("artists/%s", spotifyId))
	if err != nil {
		t.logger.LogWarn(err.Error())
		return spotifyArtists.Artist{}, nil
	}

	var spotifyArtist spotifyArtists.Artist
	err = httpClient.MakeRequest(t.httpConfig, httpRequest, &spotifyArtist)
	if err != nil {
		t.logger.LogWarn(err.Error())
	}

	return spotifyArtist, err
}

func (t *SpotifyService) GetArtists(spotifyIds []string) []spotifyArtists.Artist {
	chunkedIds := helpers.Chunk(spotifyIds, ARTIST_QUERY_LIMIT)
	httpRequests := t.getHttpRequests(chunkedIds, "GET", "artists?ids=%s")
	spotifyArtistContainers := httpClient.MakeBatchRequest[spotifyArtists.ArtistContainer](t.httpConfig, httpRequests)

	results := make([]spotifyArtists.Artist, 0)
	for _, container := range spotifyArtistContainers {
		results = append(results, container.Artists...)
	}

	return results
}

func (t *SpotifyService) GetArtistReleases(spotifyId string) []releases.Release {
	var spotifyResponse releases.ArtistReleasesContainer
	offset := 0
	for {
		offset = offset + RELEASE_QUERY_LIMIT

		httpRequest, err := t.getHttpRequest("GET", fmt.Sprintf("artists/%s/albums?limit=%d&offset=%d", spotifyId, RELEASE_QUERY_LIMIT, offset))
		if err != nil {
			t.logger.LogWarn(err.Error())
			continue
		}

		var tmpResponse releases.ArtistReleasesContainer
		err = httpClient.MakeRequest(t.httpConfig, httpRequest, &tmpResponse)
		if err != nil {
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
	chunkedIds := helpers.Chunk(spotifyIds, RELEASE_QUERY_LIMIT)
	httpRequests := t.getHttpRequests(chunkedIds, "GET", "albums?ids=%s")
	releaseContainers := httpClient.MakeBatchRequest[releases.ReleaseDetailsContainer](t.httpConfig, httpRequests)

	spotifyReleases := make([]releases.Release, 0)
	for _, container := range releaseContainers {
		spotifyReleases = append(spotifyReleases, container.Releases...)
	}

	return spotifyReleases
}

func (t *SpotifyService) SearchArtist(query string) []spotifyArtists.SlimArtist {
	httpRequest, err := t.getHttpRequest("GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)))
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]spotifyArtists.SlimArtist, 0)
	}

	var spotifyArtistSearchResults searchModels.ArtistSearchResult
	err = httpClient.MakeRequest(t.httpConfig, httpRequest, &spotifyArtistSearchResults)
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]spotifyArtists.SlimArtist, 0)
	}

	return spotifyArtistSearchResults.Artists.Items
}

func (t *SpotifyService) getHttpRequests(params [][]string, method string, format string) []*http.Request {
	httpRequests := make([]*http.Request, len(params))
	for i, param := range params {
		joinedParams := strings.Join(param, ",")
		request, err := t.getHttpRequest(method, fmt.Sprintf(format, joinedParams))
		if err != nil {
			t.logger.LogWarn("can't build an http request: %s", err.Error())
			continue
		}

		httpRequests[i] = request
	}

	return httpRequests
}

func (t *SpotifyService) getHttpRequest(method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "https://api.spotify.com/v1/"+url, nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken(t.logger, t.config)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

const ARTIST_QUERY_LIMIT = 50
const RELEASE_QUERY_LIMIT = 20
