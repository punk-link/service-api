package static

import (
	"fmt"
	artistModels "main/models/artists"
	platformModels "main/models/platforms"
	artistServices "main/services/artists"
	commonServices "main/services/common"
	platformServices "main/services/platforms"
	"time"

	cacheManager "github.com/punk-link/cache-manager"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type StaticReleaseService struct {
	cache           cacheManager.CacheManager[map[string]any]
	hashCoder       *commonServices.HashCoder
	logger          logger.Logger
	platformService *platformServices.StreamingPlatformService
	releaseService  *artistServices.ReleaseService
}

func NewStaticReleaseService(injector *do.Injector) (*StaticReleaseService, error) {
	cache := do.MustInvoke[cacheManager.CacheManager[map[string]any]](injector)
	hashCoder := do.MustInvoke[*commonServices.HashCoder](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	platformService := do.MustInvoke[*platformServices.StreamingPlatformService](injector)
	releaseService := do.MustInvoke[*artistServices.ReleaseService](injector)

	return &StaticReleaseService{
		cache:           cache,
		hashCoder:       hashCoder,
		logger:          logger,
		platformService: platformService,
		releaseService:  releaseService,
	}, nil
}

func (t *StaticReleaseService) Get(hash string) (map[string]any, error) {
	cacheKey := buildReleaseCacheKey(hash)
	value, isCached := t.cache.TryGet(cacheKey)
	if isCached {
		return value, nil
	}

	id := t.hashCoder.Decode(hash)
	release, err := t.releaseService.GetOne(id)
	tracks, err := t.buildTracks(err, release.Tracks, release.ReleaseArtists)
	platformUrls, err := t.getPlatformReleaseUrls(err, id)
	result, err := buildReleaseResult(err, release, tracks, platformUrls)

	if err == nil {
		t.cache.Set(cacheKey, result, RELEASE_CACHE_DURATION)
	}

	return result, err
}

func (t *StaticReleaseService) buildTracks(err error, tracks []artistModels.Track, releaseArtists []artistModels.Artist) ([]artistModels.SlimTrack, error) {
	if err != nil {
		return make([]artistModels.SlimTrack, 0), err
	}

	releaseArtistIds := make(map[int]int, len(releaseArtists))
	for _, artist := range releaseArtists {
		releaseArtistIds[artist.Id] = 0
	}

	slimTracks := make([]artistModels.SlimTrack, len(tracks))
	for i, track := range tracks {
		trackArtists := make([]string, 0)
		for _, artist := range track.Artists {
			if _, isExist := releaseArtistIds[artist.Id]; !isExist {
				trackArtists = append(trackArtists, artist.Name)
			}
		}

		slimTracks[i] = artistModels.SlimTrack{
			ArtistNames: trackArtists,
			IsExplicit:  track.IsExplicit,
			Name:        track.Name,
		}
	}

	return slimTracks, err
}

func (t *StaticReleaseService) getPlatformReleaseUrls(err error, id int) ([]platformModels.PlatformReleaseUrl, error) {
	if err != nil {
		return make([]platformModels.PlatformReleaseUrl, 0), err
	}

	return t.platformService.Get(id)
}

func buildReleaseCacheKey(hash string) string {
	return fmt.Sprintf("ArtistStaticRelease::%s", hash)
}

func buildReleaseResult(err error, release artistModels.Release, tracks []artistModels.SlimTrack, platformUrls []platformModels.PlatformReleaseUrl) (map[string]any, error) {
	if err != nil {
		return make(map[string]any), err
	}

	return map[string]any{
		"PageTitle":          fmt.Sprintf("%s – %s", release.Name, release.ReleaseArtists[0].Name),
		"ArtistNames":        release.ReleaseArtists,
		"ReleaseName":        release.Name,
		"ReleaseDate":        release.ReleaseDate.Year(),
		"ImageDetails":       release.ImageDetails,
		"Tracks":             tracks,
		"StreamingPlatforms": platformUrls,
	}, err
}

const RELEASE_CACHE_DURATION = time.Hour * 24
