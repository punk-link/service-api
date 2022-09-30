package spotify

import (
	"fmt"
	"main/helpers"
	commonModels "main/models/common"
	platforms "main/models/platforms"
	spotifyArtists "main/models/spotify/artists"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/services/common"
	platformServices "main/services/platforms/base"
	"net/url"
	"strings"

	"github.com/samber/do"
)

type SpotifyService struct {
	logger *common.Logger
}

func ConstructSpotifyService(injector *do.Injector) (*SpotifyService, error) {
	logger := do.MustInvoke[*common.Logger](injector)

	return &SpotifyService{
		logger: logger,
	}, nil
}

func ConstructSpotifyServiceAsPlatformer(injector *do.Injector) (platformServices.Platformer, error) {
	logger := do.MustInvoke[*common.Logger](injector)

	return platformServices.Platformer(&SpotifyService{
		logger: logger,
	}), nil
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
	var spotifyResponse releases.ArtistReleasesContainer
	offset := 0
	for {
		urlPattern := fmt.Sprintf("artists/%s/albums?limit=%d&offset=%d", spotifyId, QUERY_LIMIT, offset)
		offset = offset + QUERY_LIMIT

		var tmpResponse releases.ArtistReleasesContainer
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
	chunkedIds := helpers.Chunk(spotifyIds, QUERY_LIMIT)
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

func (t *SpotifyService) GetReleaseUrlsByUpc(upcContainers []platforms.UpcContainer) []platforms.UpcContainer {
	urlsWithSync := make([]commonModels.SyncedUrl, len(upcContainers))
	for i, container := range upcContainers {
		urlsWithSync[i] = commonModels.SyncedUrl{
			Sync: container.Upc,
			Url:  fmt.Sprintf("search?type=album&q=upc:%s", container.Upc),
		}
	}

	upcMap := make(map[string]int)
	for _, container := range upcContainers {
		upcMap[container.Upc] = container.Id
	}

	syncedReleaseContainers := makeBatchRequestWithSync[releases.UpcArtistReleasesContainer](t.logger, "GET", urlsWithSync)

	results := make([]platforms.UpcContainer, 0)
	for _, syncedContainer := range syncedReleaseContainers {
		container := syncedContainer.Result.(releases.UpcArtistReleasesContainer)

		if len(container.Albums.Items) == 0 {
			continue
		}

		release := container.Albums.Items[0]
		id := upcMap[syncedContainer.Sync]

		result := platforms.UpcContainer{
			Id:  id,
			Upc: syncedContainer.Sync,
			Url: release.ExternalUrls.Spotify,
		}

		results = append(results, result)
	}

	return results
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

const QUERY_LIMIT int = 20
