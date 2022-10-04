package static

import (
	"fmt"
	"main/models/artists"
	"main/models/platforms"
	artistServices "main/services/artists"
	"main/services/cache"
	"main/services/common"
	platformServices "main/services/platforms"
	"time"

	"github.com/samber/do"
)

type StaticReleaseService struct {
	cache           *cache.MemoryCacheService
	hashCoder       *common.HashCoder
	logger          *common.Logger
	platformService *platformServices.StrimingPlatformService
	releaseService  *artistServices.ReleaseService
}

func ConstructStaticReleaseService(injector *do.Injector) (*StaticReleaseService, error) {
	cache := do.MustInvoke[*cache.MemoryCacheService](injector)
	hashCoder := do.MustInvoke[*common.HashCoder](injector)
	logger := do.MustInvoke[*common.Logger](injector)
	platformService := do.MustInvoke[*platformServices.StrimingPlatformService](injector)
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
		return value.(map[string]any), nil
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

func (t *StaticReleaseService) buildTracks(err error, tracks []artists.Track, releaseArtists []artists.Artist) ([]artists.SlimTrack, error) {
	if err != nil {
		return make([]artists.SlimTrack, 0), err
	}

	releaseArtistIds := make(map[int]int, len(releaseArtists))
	for _, artist := range releaseArtists {
		releaseArtistIds[artist.Id] = 0
	}

	slimTracks := make([]artists.SlimTrack, len(tracks))
	for i, track := range tracks {
		trackArtists := make([]string, 0)
		for _, artist := range track.Artists {
			if _, isExist := releaseArtistIds[artist.Id]; !isExist {
				trackArtists = append(trackArtists, artist.Name)
			}
		}

		slimTracks[i] = artists.SlimTrack{
			ArtistNames: trackArtists,
			IsExplicit:  track.IsExplicit,
			Name:        track.Name,
		}
	}

	return slimTracks, err
}

func (t *StaticReleaseService) getPlatformReleaseUrls(err error, id int) ([]platforms.PlatformReleaseUrl, error) {
	if err != nil {
		return make([]platforms.PlatformReleaseUrl, 0), err
	}

	return t.platformService.Get(id)
}

func buildReleaseCacheKey(hash string) string {
	return fmt.Sprintf("ArtistStaticRelease::%s", hash)
}

func buildReleaseResult(err error, release artists.Release, tracks []artists.SlimTrack, platformUrls []platforms.PlatformReleaseUrl) (map[string]any, error) {
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
