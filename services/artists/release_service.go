package artists

import (
	"fmt"
	artistData "main/data/artists"
	artistModels "main/models/artists"
	"main/models/labels"
	"main/models/platforms"
	"main/models/platforms/spotify/releases"
	"main/services/artists/converters"
	"main/services/cache"
	"main/services/platforms/spotify"
	"sync"
	"time"

	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type ReleaseService struct {
	cache          *cache.MemoryCacheService
	logger         *logger.Logger
	spotifyService *spotify.SpotifyService
}

func ConstructReleaseService(injector *do.Injector) (*ReleaseService, error) {
	cache := do.MustInvoke[*cache.MemoryCacheService](injector)
	logger := do.MustInvoke[*logger.Logger](injector)
	spotifyService := do.MustInvoke[*spotify.SpotifyService](injector)

	return &ReleaseService{
		cache:          cache,
		logger:         logger,
		spotifyService: spotifyService,
	}, nil
}

func (t *ReleaseService) Add(currentManager labels.ManagerContext, artists map[string]artistData.Artist, releases []releases.Release, timeStamp time.Time) error {
	dbReleases := t.buildDbReleases(artists, releases, timeStamp)
	orderDbReleasesChronologically(dbReleases)

	return createDbReleasesInBatches(t.logger, nil, &dbReleases)
}

func (t *ReleaseService) GetCount() int {
	count, _ := getDbReleaseCount(t.logger, nil)
	return int(count)
}

func (t *ReleaseService) GetByArtistId(artistId int) ([]artistModels.Release, error) {
	cacheKey := t.buildArtistReleasesCacheKey(artistId)
	value, isCached := t.cache.TryGet(cacheKey)
	if isCached {
		return value.([]artistModels.Release), nil
	}

	dbReleases, err := getDbReleasesByArtistId(t.logger, nil, artistId)
	artists, err := t.getReleasesArtists(err, dbReleases)
	releases, err := t.toReleases(err, dbReleases, artists)
	if err == nil {
		t.cache.Set(cacheKey, releases, RELEASE_CACHE_DURATION)
	}

	return releases, err
}

func (t *ReleaseService) GetMissing(artistId int, artistSpotifyId string) ([]releases.Release, error) {
	dbReleases, err := getDbReleasesByArtistId(t.logger, nil, artistId)
	missingReleaseSpotifyIds, err := t.getMissingReleasesSpotifyIds(err, dbReleases, artistSpotifyId)

	return t.getReleaseDetails(err, missingReleaseSpotifyIds)
}

func (t *ReleaseService) GetOne(id int) (artistModels.Release, error) {
	cacheKey := t.buildCacheKey(id)
	value, isCached := t.cache.TryGet(cacheKey)
	if isCached {
		return value.(artistModels.Release), nil
	}

	dbRelease, err := getDbRelease(t.logger, nil, id)
	dbReleasesOfOne := []artistData.Release{dbRelease}

	artists, err := t.getReleasesArtists(err, dbReleasesOfOne)
	releases, err := t.toReleases(err, dbReleasesOfOne, artists)
	if err == nil {
		t.cache.Set(cacheKey, releases[0], RELEASE_CACHE_DURATION)
	}

	return releases[0], nil
}

func (t *ReleaseService) GetUpcContainersToUpdate(top int, skip int, updateTreshold time.Time) []platforms.UpcContainer {
	releases, _ := getUpcContainers(t.logger, nil, top, skip, updateTreshold)

	results := make([]platforms.UpcContainer, len(releases))
	for i, release := range releases {
		results[i] = platforms.UpcContainer{
			Id:  release.Id,
			Upc: release.Upc,
		}
	}

	return results
}

func (t *ReleaseService) MarkAsUpdated(ids []int, timestamp time.Time) error {
	return markDbReleasesAsUpdated(t.logger, nil, ids, timestamp)
}

func (t *ReleaseService) buildCacheKey(id int) string {
	return fmt.Sprintf("Release::%v", id)
}

func (t *ReleaseService) buildArtistReleasesCacheKey(artistId int) string {
	return fmt.Sprintf("ArtistReleases::%v", artistId)
}

func (t *ReleaseService) buildDbReleases(artists map[string]artistData.Artist, releases []releases.Release, timeStamp time.Time) []artistData.Release {
	var wg sync.WaitGroup
	chanResults := make(chan artistData.Release)
	for _, release := range releases {
		wg.Add(1)
		go t.buildFromSpotify(&wg, chanResults, artists, release, timeStamp)
	}

	go func() {
		wg.Wait()
		close(chanResults)
	}()

	dbReleases := make([]artistData.Release, 0)
	for result := range chanResults {
		dbReleases = append(dbReleases, result)
	}

	return dbReleases
}

func (t *ReleaseService) buildFromSpotify(wg *sync.WaitGroup, results chan<- artistData.Release, artists map[string]artistData.Artist, release releases.Release, timeStamp time.Time) {
	defer wg.Done()

	dbArtist, err := converters.ToDbReleaseFromSpotify(release, artists, timeStamp)
	if err != nil {
		t.logger.LogError(err, err.Error())
		return
	}

	results <- dbArtist
}

func (t *ReleaseService) getMissingReleasesSpotifyIds(err error, dbReleases []artistData.Release, artistSpotifyId string) ([]string, error) {
	if err != nil {
		return make([]string, 0), err
	}

	dbReleaseIds := make(map[string]int, len(dbReleases))
	for _, release := range dbReleases {
		dbReleaseIds[release.SpotifyId] = 0
	}

	spotifyReleases := t.spotifyService.GetArtistReleases(artistSpotifyId)
	missingReleaseSpotifyIds := make([]string, 0)
	for _, spotifyRelease := range spotifyReleases {
		if _, isContains := dbReleaseIds[spotifyRelease.Id]; !isContains {
			missingReleaseSpotifyIds = append(missingReleaseSpotifyIds, spotifyRelease.Id)
		}
	}

	return missingReleaseSpotifyIds, nil
}

func (t *ReleaseService) getReleaseDetails(err error, missingReleaseSpotifyIds []string) ([]releases.Release, error) {
	if err != nil {
		return make([]releases.Release, 0), err
	}

	details := t.spotifyService.GetReleasesDetails(missingReleaseSpotifyIds)
	return details, err
}

func (t *ReleaseService) getReleasesArtists(err error, releases []artistData.Release) (map[int]artistModels.Artist, error) {
	if err != nil {
		return make(map[int]artistModels.Artist, 0), err
	}

	artistIds := getArtistsIdsFromDbReleases(t.logger, releases)
	artists, err := getDbArtists(t.logger, err, artistIds)

	results := make(map[int]artistModels.Artist, len(artists))
	for _, dbArtist := range artists {
		artist, err := converters.ToArtist(dbArtist, make([]artistModels.Release, 0))
		if err == nil {
			results[artist.Id] = artist
		}
	}

	return results, err
}

func (t *ReleaseService) toReleases(err error, releases []artistData.Release, artists map[int]artistModels.Artist) ([]artistModels.Release, error) {
	if err != nil {
		return make([]artistModels.Release, 0), err
	}

	results, err := converters.ToReleases(releases, artists)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return results, err
}

const RELEASE_CACHE_DURATION = time.Hour * 24
