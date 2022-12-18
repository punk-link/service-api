package static

import (
	"fmt"
	artistModels "main/models/artists"
	"main/models/artists/enums"
	artistServices "main/services/artists"
	"main/services/artists/converters"
	commonServices "main/services/common"
	"sort"
	"time"

	cacheManager "github.com/punk-link/cache-manager"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type StaticArtistService struct {
	cache          cacheManager.CacheManager[map[string]any]
	hashCoder      *commonServices.HashCoder
	logger         logger.Logger
	artistService  *artistServices.ArtistService
	releaseService *artistServices.ReleaseService
}

func NewStaticArtistService(injector *do.Injector) (*StaticArtistService, error) {
	cache := do.MustInvoke[cacheManager.CacheManager[map[string]any]](injector)
	hashCoder := do.MustInvoke[*commonServices.HashCoder](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	artistService := do.MustInvoke[*artistServices.ArtistService](injector)
	releaseService := do.MustInvoke[*artistServices.ReleaseService](injector)

	return &StaticArtistService{
		cache:          cache,
		hashCoder:      hashCoder,
		logger:         logger,
		artistService:  artistService,
		releaseService: releaseService,
	}, nil
}

func (t *StaticArtistService) Get(hash string) (map[string]any, error) {
	cacheKey := buildArtistCacheKey(hash)
	value, isCached := t.cache.TryGet(cacheKey)
	if isCached {
		return value, nil
	}

	id := t.hashCoder.Decode(hash)
	artist, err := t.artistService.GetOne(id)
	releases, err := t.getReleases(err, id)
	soleReleases, compilations, err := t.sortReleases(err, releases)
	result, err := buildArtistResult(err, t.hashCoder, artist, soleReleases, compilations)

	if err == nil {
		t.cache.Set(cacheKey, result, ARTIST_CACHE_DURATION)
	}

	return result, err
}

func (t *StaticArtistService) getReleases(err error, artistId int) ([]artistModels.Release, error) {
	if err != nil {
		return make([]artistModels.Release, 0), err
	}

	return t.releaseService.GetByArtistId(artistId)
}

func (t *StaticArtistService) sortReleases(err error, releases []artistModels.Release) ([]artistModels.Release, []artistModels.Release, error) {
	if err != nil {
		return make([]artistModels.Release, 0), make([]artistModels.Release, 0), err
	}

	soleReleases := make([]artistModels.Release, 0)
	compilations := make([]artistModels.Release, 0)
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

func buildArtistCacheKey(hash string) string {
	return fmt.Sprintf("StaticArtist::%s", hash)
}

func buildArtistResult(err error, hashCoder *commonServices.HashCoder, artist artistModels.Artist, soleReleases []artistModels.Release, compilations []artistModels.Release) (map[string]any, error) {
	if err != nil {
		return make(map[string]any, 0), err
	}

	return map[string]any{
		"PageTitle":         artist.Name,
		"ArtistName":        artist.Name,
		"SoleReleaseNumber": len(soleReleases),
		"CompilationNumber": len(compilations),
		"Releases":          converters.ToSlimRelease(hashCoder, append(soleReleases, compilations...)),
	}, err
}

func sortReleasesInternal(releases []artistModels.Release) {
	sort.SliceStable(releases, func(i, j int) bool {
		return releases[i].ReleaseDate.After(releases[j].ReleaseDate)
	})
}

const ARTIST_CACHE_DURATION = time.Hour * 24
