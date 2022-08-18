package spotify

import (
	"fmt"
	"main/helpers"
	"main/models/artists"
	"main/models/spotify/releases"
	"main/models/spotify/search"
	"main/services/common"
	"net/url"
	"strings"
	"sync"
	"time"
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
	chunkedIds := helpers.Chunk(spotifyIds, queryLimit)
	mainLoop, extraLoop := t.divideChunkToLoops(chunkedIds, iterationStep)

	var wg sync.WaitGroup
	results := make(chan []releases.Release)
	for i := 0; i < len(mainLoop); i = i + iterationStep {
		wg.Add(iterationStep)

		for j := 0; j < iterationStep; j++ {
			idQuery := strings.Join(mainLoop[i+j], ",")
			go t.getReleasesDetails(idQuery, &wg, results)
		}

		time.Sleep(requestBatchTimeoutDuration)
	}

	for _, chunk := range extraLoop {
		wg.Add(1)
		idQuery := strings.Join(chunk, ",")
		go t.getReleasesDetails(idQuery, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	spotifyReleases := make([]releases.Release, 0)
	for result := range results {
		spotifyReleases = append(spotifyReleases, result...)
	}

	return spotifyReleases
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

func (t *SpotifyService) SearchArtist(query string) []search.Artist {
	var spotifyArtistSearchResults search.ArtistSearchResult
	err := makeRequest(t.logger, "GET", fmt.Sprintf("search?type=artist&limit=10&q=%s", url.QueryEscape(query)), &spotifyArtistSearchResults)
	if err != nil {
		t.logger.LogWarn(err.Error())
		return make([]search.Artist, 0)
	}

	return spotifyArtistSearchResults.Artists.Items
}

func (t *SpotifyService) divideChunkToLoops(chunkedIds [][]string, iterationStep int) ([][]string, [][]string) {
	var extraLoop [][]string
	mainLoop := make([][]string, 0)

	if len(chunkedIds) < iterationStep {
		extraLoop = chunkedIds
	} else {
		extraElements := len(chunkedIds) % iterationStep

		mainLoop = chunkedIds[0 : len(chunkedIds)-extraElements]
		extraLoop = chunkedIds[len(chunkedIds)-extraElements:]
	}

	return mainLoop, extraLoop
}

func (t *SpotifyService) getReleasesDetails(idQuery string, wg *sync.WaitGroup, results chan<- []releases.Release) {
	defer wg.Done()

	var tmpResponse releases.ReleaseDetailsContainer
	if err := makeRequest(t.logger, "GET", fmt.Sprintf("albums?ids=%s", idQuery), &tmpResponse); err != nil {
		t.logger.LogWarn(err.Error())
		return
	}

	results <- tmpResponse.Items
}

const queryLimit int = 20
const iterationStep int = 4
const requestBatchTimeoutDuration time.Duration = time.Millisecond * 100
