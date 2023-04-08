package spotify

import (
	"fmt"
	"main/helpers"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type SpotifyReleaseService struct {
	httpConfig    *httpClient.HttpClientConfig
	logger        logger.Logger
	spotifyClient SpotifyClient
}

func NewSpotifyReleaseService(injector *do.Injector) (SpotifyReleaseServer, error) {
	httpConfig := do.MustInvoke[*httpClient.HttpClientConfig](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	spotifyClient := do.MustInvoke[SpotifyClient](injector)

	return &SpotifyReleaseService{
		httpConfig:    httpConfig,
		logger:        logger,
		spotifyClient: spotifyClient,
	}, nil
}

func (t *SpotifyReleaseService) Get(spotifyIds []string) []releaseSpotifyPlatformModels.Release {
	chunkedIds := helpers.Chunk(spotifyIds, RELEASE_QUERY_LIMIT)
	httpRequests := t.spotifyClient.Request(chunkedIds, "GET", "albums?ids=%s")

	httpClient := httpClient.New[releaseSpotifyPlatformModels.ReleaseDetailsContainer](t.httpConfig)
	releaseContainers := httpClient.MakeBatchRequest(httpRequests)

	spotifyReleases := make([]releaseSpotifyPlatformModels.Release, 0)
	for _, container := range releaseContainers {
		spotifyReleases = append(spotifyReleases, container.Releases...)
	}

	return spotifyReleases
}

func (t *SpotifyReleaseService) GetByArtistId(spotifyId string) []releaseSpotifyPlatformModels.Release {
	var spotifyResponse releaseSpotifyPlatformModels.ArtistReleasesContainer
	offset := 0
	for {
		offset = offset + RELEASE_QUERY_LIMIT

		httpRequest, err := t.spotifyClient.RequestOne("GET", fmt.Sprintf("artists/%s/albums?limit=%d&offset=%d", spotifyId, RELEASE_QUERY_LIMIT, offset))
		if err != nil {
			t.logger.LogWarn(err.Error())
			continue
		}

		httpClient := httpClient.New[releaseSpotifyPlatformModels.ArtistReleasesContainer](t.httpConfig)
		tmpResponse, err := httpClient.MakeRequest(httpRequest)
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

const RELEASE_QUERY_LIMIT = 20
