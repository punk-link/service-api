package artists

import (
	"fmt"
	"main/models/artists"
	"main/services/cache"
	"main/services/common"
	"strings"
)

type MvcReleaseService struct {
	cache          *cache.MemoryCacheService
	hashCoder      *common.HashCoder
	logger         *common.Logger
	releaseService *ReleaseService
}

func ConstructMvcReleaseService(cache *cache.MemoryCacheService, logger *common.Logger, releaseService *ReleaseService, hashCoder *common.HashCoder) *MvcReleaseService {
	return &MvcReleaseService{
		cache:          cache,
		hashCoder:      hashCoder,
		logger:         logger,
		releaseService: releaseService,
	}
}

func (t *MvcReleaseService) Get(hash string) (map[string]any, error) {
	cacheKey := t.buildReleaseCacheKey(hash)
	value, isCached := t.cache.TryGet(cacheKey)
	if isCached {
		return value.(map[string]any), nil
	}

	id := t.hashCoder.Encode(hash)
	release, err := t.releaseService.GetOne(id)
	artistNames, err := t.buildArtistNames(err, release.ReleaseArtists)
	tracks, err := t.buildTracks(err, release.Tracks)
	if err != nil {
		return make(map[string]any), err
	}

	result := map[string]any{
		"PageTitle":         fmt.Sprintf("%s – %s", release.Name, release.ReleaseArtists[0].Name),
		"ArtistName":        artistNames,
		"ReleaseName":       release.Name,
		"ReleaseDate":       release.ReleaseDate.Year(),
		"ImageDetails":      release.ImageDetails,
		"Tracks":            tracks,
		"StreamingServices": []string{"Apple Music", "Deezer"},
	}

	t.cache.Set(cacheKey, result, RELEASE_CACHE_DURATION)

	return result, err
}

func (t *MvcReleaseService) buildReleaseCacheKey(hash string) string {
	return fmt.Sprintf("ArtistMvcRelease::%s", hash)
}

func (t *MvcReleaseService) buildArtistNames(err error, artists []artists.Artist) (string, error) {
	if err != nil {
		return "", err
	}

	names := make([]string, len(artists))
	for i, artist := range artists {
		names[i] = artist.Name
	}

	return strings.Join(names, ", "), err
}

func (t *MvcReleaseService) buildTracks(err error, tracks []artists.Track) ([]artists.SlimTrack, error) {
	if err != nil {
		return make([]artists.SlimTrack, 0), err
	}

	slimTracks := make([]artists.SlimTrack, len(tracks))
	for i, track := range tracks {
		slimTracks[i] = artists.SlimTrack{
			ArtistNames: track.Artists[0].Name,
			IsExplicit:  track.IsExplicit,
			Name:        track.Name,
		}
	}

	return slimTracks, err
}
