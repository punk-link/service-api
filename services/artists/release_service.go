package artists

import (
	dataStructures "main/data-structures"
	artistData "main/data/artists"
	"main/helpers"
	artistModels "main/models/artists"
	labelModels "main/models/labels"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"
	"main/services/artists/converters"
	"main/services/artists/extractors"
	"main/services/artists/repositories"
	spotifyPlatformServices "main/services/platforms/spotify"
	"sort"
	"sync"
	"time"

	cacheManager "github.com/punk-link/cache-manager"
	"github.com/punk-link/logger"
	platformContracts "github.com/punk-link/platform-contracts"
	"github.com/samber/do"
)

type ReleaseService struct {
	artistIdExtractor     extractors.ArtistIdExtractor
	artistRepository      repositories.ArtistRepository
	logger                logger.Logger
	releaseCache          cacheManager.CacheManager[artistModels.Release]
	releasesCache         cacheManager.CacheManager[[]artistModels.Release]
	repository            repositories.ReleaseRepository
	spotifyReleaseService spotifyPlatformServices.SpotifyReleaseServer
	tagService            TagServer
}

func NewReleaseService(injector *do.Injector) (ReleaseServer, error) {
	artistIdExtractor := do.MustInvoke[extractors.ArtistIdExtractor](injector)
	artistRepository := do.MustInvoke[repositories.ArtistRepository](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	releaseCache := do.MustInvoke[cacheManager.CacheManager[artistModels.Release]](injector)
	releasesCache := do.MustInvoke[cacheManager.CacheManager[[]artistModels.Release]](injector)
	repository := do.MustInvoke[repositories.ReleaseRepository](injector)
	spotifyReleaseService := do.MustInvoke[spotifyPlatformServices.SpotifyReleaseServer](injector)
	tagService := do.MustInvoke[TagServer](injector)

	return &ReleaseService{
		artistIdExtractor:     artistIdExtractor,
		artistRepository:      artistRepository,
		logger:                logger,
		releaseCache:          releaseCache,
		releasesCache:         releasesCache,
		repository:            repository,
		spotifyReleaseService: spotifyReleaseService,
		tagService:            tagService,
	}, nil
}

func (t *ReleaseService) Add(currentManager labelModels.ManagerContext, artists map[string]artistData.Artist, releases []releaseSpotifyPlatformModels.Release, timeStamp time.Time) error {
	dbReleases := t.buildDbReleases(artists, releases, timeStamp)
	t.orderDbReleasesChronologically(dbReleases)

	err := t.repository.CreateInBatches(nil, &dbReleases)
	if err != nil {
		return err
	}

	spotifyReleaseTags := make(map[string][]string, len(releases))
	for _, release := range releases {
		spotifyReleaseTags[release.Id] = release.Genres
	}

	releaseTagRelations := make([]artistData.ReleaseTagRelation, 0)
	for _, dbRelease := range dbReleases {
		tagSlugs := spotifyReleaseTags[dbRelease.SpotifyId]
		tags := t.tagService.GetOrAdd(tagSlugs)
		for _, tag := range tags {
			releaseTagRelations = append(releaseTagRelations, converters.ToDbReleaseTagRelation(dbRelease.Id, tag.Id))
		}
	}

	err = t.repository.AddTags(err, &releaseTagRelations)

	return err
}

func (t *ReleaseService) Get(artistId int) ([]artistModels.Release, error) {
	cacheKey := t.buildArtistReleasesCacheKey(artistId)
	value, isCached := t.releasesCache.TryGet(cacheKey)
	if isCached {
		return value, nil
	}

	dbReleases, err := t.repository.Get(nil, artistId)
	releases, err := t.getReleasesInternal(err, dbReleases)
	if err == nil {
		t.releasesCache.Set(cacheKey, releases, RELEASE_CACHE_DURATION)
	}

	return releases, err
}

func (t *ReleaseService) GetCount() int {
	count, _ := t.repository.GetCount(nil)
	return int(count)
}

func (t *ReleaseService) GetMissing(artistId int, artistSpotifyId string) ([]releaseSpotifyPlatformModels.Release, error) {
	dbReleases, err := t.repository.Get(nil, artistId)
	missingReleaseSpotifyIds, err := t.getMissingReleasesSpotifyIds(err, dbReleases, artistSpotifyId)

	return t.getReleaseDetails(err, missingReleaseSpotifyIds)
}

func (t *ReleaseService) GetOne(id int) (artistModels.Release, error) {
	cacheKey := t.buildCacheKey(id)
	value, isCached := t.releaseCache.TryGet(cacheKey)
	if isCached {
		return value, nil
	}

	dbRelease, err := t.repository.GetOne(nil, id)
	dbReleasesOfOne := []artistData.Release{dbRelease}
	releases, err := t.getReleasesInternal(err, dbReleasesOfOne)
	if err == nil {
		t.releaseCache.Set(cacheKey, releases[0], RELEASE_CACHE_DURATION)
	}

	return releases[0], nil
}

func (t *ReleaseService) GetUpcContainersToUpdate(top int, skip int, updateTreshold time.Time) []platformContracts.UpcContainer {
	releases, _ := t.repository.GetUpcContainers(nil, top, skip, updateTreshold)

	results := make([]platformContracts.UpcContainer, len(releases))
	for i, release := range releases {
		results[i] = platformContracts.UpcContainer{
			Id:  release.Id,
			Upc: release.Upc,
		}
	}

	return results
}

func (t *ReleaseService) MarkAsUpdated(ids []int, timestamp time.Time) error {
	return t.repository.MarksAsUpdated(nil, ids, timestamp)
}

func (t *ReleaseService) buildCacheKey(id int) string {
	return helpers.BuildCacheKey("Release", id)
}

func (t *ReleaseService) buildArtistReleasesCacheKey(artistId int) string {
	return helpers.BuildCacheKey("ArtistReleases", artistId)
}

func (t *ReleaseService) buildDbReleases(artists map[string]artistData.Artist, releases []releaseSpotifyPlatformModels.Release, timeStamp time.Time) []artistData.Release {
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

func (t *ReleaseService) buildFromSpotify(wg *sync.WaitGroup, results chan<- artistData.Release, artists map[string]artistData.Artist, release releaseSpotifyPlatformModels.Release, timeStamp time.Time) {
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

	dbReleaseIds := dataStructures.MakeEmptyHashSet[string]()
	for _, release := range dbReleases {
		dbReleaseIds.Add(release.SpotifyId)
	}

	spotifyReleases := t.spotifyReleaseService.GetByArtistId(artistSpotifyId)
	missingReleaseSpotifyIds := make([]string, 0)
	for _, spotifyRelease := range spotifyReleases {
		if !dbReleaseIds.Contains(spotifyRelease.Id) {
			missingReleaseSpotifyIds = append(missingReleaseSpotifyIds, spotifyRelease.Id)
		}
	}

	return missingReleaseSpotifyIds, nil
}

func (t *ReleaseService) getReleaseDetails(err error, missingReleaseSpotifyIds []string) ([]releaseSpotifyPlatformModels.Release, error) {
	if err != nil {
		return make([]releaseSpotifyPlatformModels.Release, 0), err
	}

	details := t.spotifyReleaseService.Get(missingReleaseSpotifyIds)
	return details, err
}

func (t *ReleaseService) getReleasesInternal(err error, dbReleases []artistData.Release) ([]artistModels.Release, error) {
	if err != nil {
		return make([]artistModels.Release, 0), err
	}

	artists, err := t.getReleasesArtists(err, dbReleases)
	tags, err := t.getReleasesTags(err, dbReleases)
	return t.toReleases(err, dbReleases, artists, tags)
}

func (t *ReleaseService) getReleasesArtists(err error, releases []artistData.Release) (map[int]artistModels.Artist, error) {
	if err != nil {
		return make(map[int]artistModels.Artist, 0), err
	}

	artistIds := t.artistIdExtractor.Extract(releases)
	artists, err := t.artistRepository.Get(err, artistIds)

	results := make(map[int]artistModels.Artist, len(artists))
	for _, dbArtist := range artists {
		artist, err := converters.ToArtist(dbArtist, make([]artistModels.Release, 0))
		if err == nil {
			results[artist.Id] = artist
		}
	}

	return results, err
}

func (t *ReleaseService) getReleasesTags(err error, releases []artistData.Release) (map[int][]artistModels.Tag, error) {
	if err != nil {
		return make(map[int][]artistModels.Tag, 0), err
	}

	releaseIds := make([]int, len(releases))
	for i, release := range releases {
		releaseIds[i] = release.Id
	}

	dbTags, err := t.repository.GetTags(err, releaseIds)
	if err != nil {
		return make(map[int][]artistModels.Tag, 0), err
	}

	tags := make(map[int][]artistModels.Tag, len(dbTags))
	for key, dbTagArray := range dbTags {
		releaseTags := make([]artistModels.Tag, len(dbTagArray))
		for i, dbTag := range dbTagArray {
			releaseTags[i] = artistModels.Tag{
				Id:             dbTag.Id,
				Name:           dbTag.Name,
				NormalizedName: dbTag.NormalizedName,
			}
		}

		tags[key] = releaseTags
	}

	return tags, nil
}

func (t *ReleaseService) orderDbReleasesChronologically(target []artistData.Release) {
	sort.Slice(target, func(i, j int) bool {
		return target[i].ReleaseDate.Before(target[j].ReleaseDate)
	})
}

func (t *ReleaseService) toReleases(err error, releases []artistData.Release, artists map[int]artistModels.Artist, tags map[int][]artistModels.Tag) ([]artistModels.Release, error) {
	if err != nil {
		return make([]artistModels.Release, 0), err
	}

	results := converters.ToReleases(t.logger, releases, artists, tags)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return results, err
}

const RELEASE_CACHE_DURATION = time.Hour * 24
