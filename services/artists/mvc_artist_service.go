package artists

import (
	"fmt"
	"main/models/artists"
	"main/models/artists/enums"
	"main/services/artists/converters"
	"main/services/cache"
	"main/services/common"
	"sort"
)

type MvcArtistService struct {
	cache          *cache.MemoryCacheService
	hashCoder      *common.HashCoder
	logger         *common.Logger
	artistService  *ArtistService
	releaseService *ReleaseService
}

func ConstructMvcArtistService(cache *cache.MemoryCacheService, logger *common.Logger, artistService *ArtistService, releaseService *ReleaseService, hashCoder *common.HashCoder) *MvcArtistService {
	return &MvcArtistService{
		cache:          cache,
		hashCoder:      hashCoder,
		logger:         logger,
		artistService:  artistService,
		releaseService: releaseService,
	}
}

func (t *MvcArtistService) Get(hash string) (map[string]any, error) {
	cacheKey := t.buildArtistCacheKey(hash)
	value, isCached := t.cache.TryGet(cacheKey)
	if isCached {
		return value.(map[string]any), nil
	}

	id := t.hashCoder.Decode(hash)
	artist, err := t.artistService.GetOne(id)
	releases, err := t.getReleases(err, id)
	soleReleases, compilations, err := t.sortReleases(err, releases)

	result := map[string]any{
		"PageTitle":         artist.Name,
		"ArtistName":        artist.Name,
		"SoleReleaseNumber": len(soleReleases),
		"CompilationNumber": len(compilations),
		"Releases":          converters.ToSlimRelease(t.hashCoder, append(soleReleases, compilations...)),
	}

	//t.cache.Set(cacheKey, result, RELEASE_CACHE_DURATION)

	return result, err
}

func (t *MvcArtistService) buildArtistCacheKey(hash string) string {
	return fmt.Sprintf("MvcArtist::%s", hash)
}

func (t *MvcArtistService) getReleases(err error, artistId int) ([]artists.Release, error) {
	if err != nil {
		return make([]artists.Release, 0), err
	}

	return t.releaseService.GetByArtistId(artistId)
}

func (t *MvcArtistService) sortReleases(err error, releases []artists.Release) ([]artists.Release, []artists.Release, error) {
	if err != nil {
		return make([]artists.Release, 0), make([]artists.Release, 0), err
	}

	soleReleases := make([]artists.Release, 0)
	compilations := make([]artists.Release, 0)
	for _, release := range releases {
		if release.Type == enums.Compilation {
			compilations = append(compilations, release)
		} else {
			soleReleases = append(soleReleases, release)
		}
	}

	sortReleasesInternal(soleReleases)
	sortReleasesInternal(compilations)

	return soleReleases, compilations, err
}

func sortReleasesInternal(releases []artists.Release) {
	sort.SliceStable(releases, func(i, j int) bool {
		return releases[i].ReleaseDate.After(releases[j].ReleaseDate)
	})
}
