package artists

import (
	"errors"
	artistData "main/data/artists"
	"main/helpers"
	artistModels "main/models/artists"
	labelModels "main/models/labels"
	releaseSpotifyPlatformModels "main/models/platforms/spotify/releases"
	"main/services/artists/converters"
	"main/services/artists/validators"
	spotifyPlatformServices "main/services/platforms/spotify"
	"time"

	cacheManager "github.com/punk-link/cache-manager"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type ArtistService struct {
	cache                cacheManager.CacheManager[artistModels.Artist]
	logger               logger.Logger
	releaseService       ReleaseServer
	repository           ArtistRepository
	spotifyArtistService spotifyPlatformServices.SpotifyArtistServer
}

func NewArtistService(injector *do.Injector) (ArtistServer, error) {
	cache := do.MustInvoke[cacheManager.CacheManager[artistModels.Artist]](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	releaseService := do.MustInvoke[ReleaseServer](injector)
	repository := do.MustInvoke[ArtistRepository](injector)
	spotifyArtistService := do.MustInvoke[spotifyPlatformServices.SpotifyArtistServer](injector)

	return &ArtistService{
		cache:                cache,
		logger:               logger,
		releaseService:       releaseService,
		repository:           repository,
		spotifyArtistService: spotifyArtistService,
	}, nil
}

func (t *ArtistService) Add(currentManager labelModels.ManagerContext, spotifyId string) (artistModels.Artist, error) {
	var err error
	if spotifyId == "" {
		err = errors.New("artist's spotify ID is empty")
	}

	dbArtist, err := t.repository.GetOneBySpotifyId(err, spotifyId)

	now := time.Now().UTC()
	if dbArtist.Id != 0 {
		err = validators.CurrentDbArtistBelongsToLabel(err, dbArtist, currentManager.LabelId)
		dbArtist, err = t.updateLabelIfNeeded(err, dbArtist, currentManager.LabelId)
	} else {
		dbArtist, err = t.addArtist(err, spotifyId, currentManager.LabelId, now)
	}

	err = t.findAndAddMissingReleases(err, currentManager, dbArtist, now)
	artist, err := t.getInternal(err, []int{dbArtist.Id})
	if err != nil {
		return artistModels.Artist{}, err
	}

	return artist[0], nil
}

func (t *ArtistService) Get(labelId int) ([]artistModels.Artist, error) {
	dbArtistIds, err := t.repository.GetIdsByLabelId(nil, labelId)
	return t.getInternal(err, dbArtistIds)
}

func (t *ArtistService) GetOne(id int) (artistModels.Artist, error) {
	artists, err := t.getInternal(nil, []int{id})
	if err != nil {
		return artistModels.Artist{}, err
	}

	return artists[0], nil
}

func (t *ArtistService) Search(query string) []artistModels.ArtistSearchResult {
	const minimalQueryLength = 3
	if len(query) < minimalQueryLength {
		return make([]artistModels.ArtistSearchResult, 0)
	}

	searchResults := t.spotifyArtistService.Search(query)
	return converters.ToArtistSearchResults(searchResults)
}

func (t *ArtistService) addArtist(err error, spotifyId string, labelId int, timeStamp time.Time) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	spotifyArtist, err := t.spotifyArtistService.GetOne(spotifyId)
	if err != nil {
		return artistData.Artist{}, err
	}

	dbArtist, err := converters.ToDbArtist(spotifyArtist, labelId, timeStamp)
	err = t.repository.Create(err, &dbArtist)

	return dbArtist, err
}

func (t *ArtistService) addMissingArtists(spotifyIds []string, timeStamp time.Time) ([]artistData.Artist, error) {
	missingArtists := t.spotifyArtistService.Get(spotifyIds)

	var err error
	dbArtists := make([]artistData.Artist, len(missingArtists))
	for i, artist := range missingArtists {
		dbArtist, localError := converters.ToDbArtist(&artist, 0, timeStamp)
		if localError != nil {
			err = helpers.CombineErrors(err, localError)
			continue
		}

		dbArtists[i] = dbArtist
	}

	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	err = t.repository.CreateInBatches(err, &dbArtists)
	return dbArtists, err
}

func (t *ArtistService) addMissingFeaturingArtists(err error, spotifyIds []string, timeStamp time.Time) (map[string]artistData.Artist, error) {
	if err != nil {
		return make(map[string]artistData.Artist, 0), err
	}

	results := make(map[string]artistData.Artist)
	addedDbArtists, err := t.addMissingArtists(spotifyIds, timeStamp)
	if err == nil {
		for _, dbArtist := range addedDbArtists {
			results[dbArtist.SpotifyId] = dbArtist
		}
	}

	return results, err
}

func (t *ArtistService) getExistingFeaturingArtists(err error, dbArtist artistData.Artist, spotifyIds []string, timeStamp time.Time) (map[string]artistData.Artist, error) {
	if err != nil {
		return make(map[string]artistData.Artist, 0), err
	}

	existedArtists, err := t.repository.GetBySpotifyIds(err, spotifyIds)

	results := make(map[string]artistData.Artist, 0)
	for _, artist := range existedArtists {
		results[artist.SpotifyId] = artist
	}

	return results, err
}

func (t *ArtistService) findAndAddMissingReleases(err error, currentManager labelModels.ManagerContext, dbArtist artistData.Artist, timeStamp time.Time) error {
	if err != nil {
		return err
	}

	missingReleases, err := t.releaseService.GetMissing(dbArtist.Id, dbArtist.SpotifyId)
	artistSpotifyIds, err := t.getFeaturingArtistSpotifyIds(err, missingReleases)
	artists, err := t.getExistingFeaturingArtists(err, dbArtist, artistSpotifyIds, timeStamp)
	missingFeaturingArtistsSpotifyIds, err := t.getMissingFeaturingArtistsSpotifyIds(err, artists, artistSpotifyIds)
	addedArtists, err := t.addMissingFeaturingArtists(err, missingFeaturingArtistsSpotifyIds, timeStamp)

	if err != nil {
		return err
	}

	for key, artist := range addedArtists {
		artists[key] = artist
	}

	return t.releaseService.Add(currentManager, artists, missingReleases, timeStamp)
}

func (t *ArtistService) getFeaturingArtistSpotifyIds(err error, releases []releaseSpotifyPlatformModels.Release) ([]string, error) {
	if err != nil {
		return make([]string, 0), err
	}

	artistIds := make(map[string]int)
	for _, release := range releases {
		for _, artist := range release.Artists {
			if _, isExists := artistIds[artist.Id]; !isExists {
				artistIds[artist.Id] = 0
			}
		}

		for _, track := range release.Tracks.Items {
			for _, artist := range track.Artists {
				if _, isExists := artistIds[artist.Id]; !isExists {
					artistIds[artist.Id] = 0
				}
			}
		}
	}

	spotifyIds := make([]string, 0)
	for i := range artistIds {
		spotifyIds = append(spotifyIds, i)
	}

	return spotifyIds, nil
}

func (t *ArtistService) getInternal(err error, ids []int) ([]artistModels.Artist, error) {
	if err != nil {
		return make([]artistModels.Artist, 0), err
	}

	artists := make([]artistModels.Artist, len(ids))
	for i, id := range ids {
		cacheKey := helpers.BuildCacheKey(ARTIST_CACHE_SLUG, id)
		value, isCached := t.cache.TryGet(cacheKey)
		if isCached {
			artists[i] = value
			continue
		}

		dbArtist, dbArtistErr := t.repository.GetOne(err, id)
		artist, conversionErr := converters.ToArtist(dbArtist, make([]artistModels.Release, 0))
		if err != nil {
			t.logger.LogError(err, err.Error())
			err = helpers.CombineErrors(err, helpers.AccumulateErrors(dbArtistErr, conversionErr))
			continue
		}

		t.cache.Set(cacheKey, artist, ARTIST_CACHE_DURATION)

		artists[i] = artist
	}

	return artists, err
}

func (t *ArtistService) getMissingFeaturingArtistsSpotifyIds(err error, existingArtists map[string]artistData.Artist, artistSpotifyIds []string) ([]string, error) {
	if err != nil {
		return make([]string, 0), err
	}

	existingArtistSpotifyIds := make(map[string]int, len(existingArtists))
	for _, artist := range existingArtists {
		existingArtistSpotifyIds[artist.SpotifyId] = 0
	}

	missingSpotifyIds := make([]string, 0)
	for _, id := range artistSpotifyIds {
		if _, isExists := existingArtistSpotifyIds[id]; !isExists {
			missingSpotifyIds = append(missingSpotifyIds, id)
		}
	}

	return missingSpotifyIds, nil
}

func (t *ArtistService) updateLabelIfNeeded(err error, dbArtist artistData.Artist, labelId int) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	if dbArtist.LabelId == 0 {
		dbArtist.LabelId = labelId
		err = t.repository.Update(err, &dbArtist)
	}

	t.cache.Remove(helpers.BuildCacheKey(ARTIST_CACHE_SLUG, dbArtist.Id))
	return dbArtist, err
}

const ARTIST_CACHE_DURATION = time.Hour * 24
const ARTIST_CACHE_SLUG = "Artist"
