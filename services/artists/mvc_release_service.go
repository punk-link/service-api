package artists

import (
	"errors"
	"fmt"
	"main/models/artists"
	"main/services/cache"
	"main/services/common"
	"strconv"
	"strings"
)

type MvcReleaseService struct {
	cache          *cache.MemoryCacheService
	logger         *common.Logger
	releaseService *ReleaseService
}

func ConstructMvcReleaseService(cache *cache.MemoryCacheService, logger *common.Logger, releaseService *ReleaseService) *MvcReleaseService {
	return &MvcReleaseService{
		cache:          cache,
		logger:         logger,
		releaseService: releaseService,
	}
}

func (t *MvcReleaseService) Get(hash string) (map[string]any, error) {
	id, err := t.getIdFromHash(hash)
	release, err := t.getRelease(err, id)
	artistNames, err := t.buildArtistNames(err, release.ReleaseArtists)
	tracks, err := t.buildTracks(err, release.Tracks)
	if err != nil {
		return make(map[string]any), err
	}

	return map[string]any{
		"PageTitle":    fmt.Sprintf("%s – %s", release.Name, release.ReleaseArtists[0].Name),
		"ArtistName":   artistNames,
		"ReleaseName":  release.Name,
		"ReleaseDate":  release.ReleaseDate.Year(),
		"ImageDetails": release.ImageDetails,
		"Tracks":       tracks,
		"Services":     []string{"Apple Music", "Deezer"},
	}, err
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

func (t *MvcReleaseService) getIdFromHash(hash string) (int, error) {
	if hash == "" {
		return 0, errors.New("")
	}

	return strconv.Atoi(hash)
}

func (t *MvcReleaseService) getRelease(err error, id int) (artists.Release, error) {
	if err != nil {
		return artists.Release{}, err
	}

	return t.releaseService.GetOne(id)
}
