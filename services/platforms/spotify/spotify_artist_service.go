package spotify

import (
	"fmt"
	"main/helpers"
	artistSpotifyPlatformModels "main/models/platforms/spotify/artists"
	searchSpotifyPlatformModels "main/models/platforms/spotify/searches"
	tokenSpotifyPlatformModels "main/models/platforms/spotify/tokens"
	"net/url"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type SpotifyArtistService struct {
	config        *tokenSpotifyPlatformModels.SpotifyClientConfig
	httpConfig    *httpClient.HttpClientConfig
	logger        logger.Logger
	spotifyClient SpotifyClient
}

func NewSpotifyArtistService(injector *do.Injector) (SpotifyArtistServer, error) {
	config := do.MustInvoke[*tokenSpotifyPlatformModels.SpotifyClientConfig](injector)
	httpConfig := do.MustInvoke[*httpClient.HttpClientConfig](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	spotifyClient := do.MustInvoke[SpotifyClient](injector)

	return &SpotifyArtistService{
		config:        config,
		httpConfig:    httpConfig,
		logger:        logger,
		spotifyClient: spotifyClient,
	}, nil
}

func (t *SpotifyArtistService) Get(spotifyIds []string) []artistSpotifyPlatformModels.Artist {
	chunkedIds := helpers.Chunk(spotifyIds, ARTIST_QUERY_LIMIT)
	httpRequests := t.spotifyClient.Request(chunkedIds, "GET", "artists?ids=%s")

	httpClient := httpClient.New[artistSpotifyPlatformModels.ArtistContainer](t.httpConfig)
	spotifyArtistContainers := httpClient.MakeBatchRequest(httpRequests)

	results := make([]artistSpotifyPlatformModels.Artist, 0)
	for _, container := range spotifyArtistContainers {
		results = append(results, container.Artists...)
	}

	return results
}

func (t *SpotifyArtistService) GetOne(spotifyId string) (*artistSpotifyPlatformModels.Artist, error) {
	httpRequest, err := t.spotifyClient.RequestOne("GET", fmt.Sprintf("artists/%s", spotifyId))
	if err != nil {
		t.logger.LogWarn(err.Error())
		return nil, err
	}

	httpClient := httpClient.New[artistSpotifyPlatformModels.Artist](t.httpConfig)
	spotifyArtist, err := httpClient.MakeRequest(httpRequest)
	if err != nil {
		t.logger.LogWarn(err.Error())
	}

	return spotifyArtist, err
}

func (t *SpotifyArtistService) Search(query string) []artistSpotifyPlatformModels.SlimArtist {
	httpRequest, err := t.spotifyClient.RequestOne("GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)))
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]artistSpotifyPlatformModels.SlimArtist, 0)
	}

	httpClient := httpClient.New[searchSpotifyPlatformModels.ArtistSearchResult](t.httpConfig)
	spotifyArtistSearchResults, err := httpClient.MakeRequest(httpRequest)
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]artistSpotifyPlatformModels.SlimArtist, 0)
	}

	return spotifyArtistSearchResults.Artists.Items
}

const ARTIST_QUERY_LIMIT = 50
