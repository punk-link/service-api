package artists

import (
	"errors"
	"fmt"
	"main/data"
	artistData "main/data/artists"
	"main/helpers"
	artistModels "main/models/artists"
	"main/models/labels"
	"main/models/spotify/releases"
	"main/services/artists/converters"
	"main/services/artists/validators"
	"main/services/common"
	"main/services/spotify"
	"time"
)

type ArtistService struct {
	cache          *ArtistCacheService
	logger         *common.Logger
	spotifyService *spotify.SpotifyService
	releaseService *ReleaseService
}

func ConstructArtistService(cache *ArtistCacheService, logger *common.Logger, releaseService *ReleaseService, spotifyService *spotify.SpotifyService) *ArtistService {
	return &ArtistService{
		cache:          cache,
		logger:         logger,
		releaseService: releaseService,
		spotifyService: spotifyService,
	}
}

func (t *ArtistService) Add(currentManager labels.ManagerContext, spotifyId string) (artistModels.Artist, error) {
	var err error
	if spotifyId == "" {
		err = errors.New("artist's spotify ID is empty")
	}

	dbArtist, err := t.getDbArtistBySpotifyId(err, spotifyId)

	now := time.Now().UTC()
	if dbArtist.Id != 0 {
		err = validators.CurrentDbArtistBelongsToLabel(err, dbArtist, currentManager.LabelId)
		dbArtist, err = t.updateLabelIfNeeded(err, dbArtist, currentManager.LabelId)
	} else {
		dbArtist, err = t.addArtist(err, spotifyId, currentManager.LabelId, now)
	}

	err = t.findAndAddMissingReleases(err, currentManager, dbArtist, now)
	return t.getInternal(err, currentManager, dbArtist.Id)
}

func (t *ArtistService) GetOne(currentManager labels.ManagerContext, id int) (artistModels.Artist, error) {
	return t.getInternal(nil, currentManager, id)
}

func (t *ArtistService) Search(query string) []artistModels.ArtistSearchResult {
	const minimalQueryLength int = 3
	if len(query) < minimalQueryLength {
		return make([]artistModels.ArtistSearchResult, 0)
	}

	searchResults := t.spotifyService.SearchArtist(query)
	return converters.ToArtistSearchResults(searchResults)
}

func (t *ArtistService) addArtist(err error, spotifyId string, labelId int, timeStamp time.Time) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	spotifyArtist, err := t.spotifyService.GetArtist(spotifyId)
	if err != nil {
		return artistData.Artist{}, err
	}

	dbArtist, err := converters.ToDbArtist(spotifyArtist, labelId, timeStamp)
	if err != nil {
		return artistData.Artist{}, err
	}

	err = data.DB.Create(&dbArtist).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return dbArtist, err
}

func (t *ArtistService) addMissingArtists(spotifyIds []string, timeStamp time.Time) ([]artistData.Artist, error) {
	missingArtists := t.spotifyService.GetArtists(spotifyIds)

	var err error
	dbArtists := make([]artistData.Artist, len(missingArtists))
	for i, artist := range missingArtists {
		dbArtist, localError := converters.ToDbArtist(artist, 0, timeStamp)
		if localError != nil {
			err = helpers.CombineErrors(err, localError)
			continue
		}

		dbArtists[i] = dbArtist
	}

	if err != nil {
		return make([]artistData.Artist, 0), err
	}

	err = data.DB.CreateInBatches(&dbArtists, 50).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return dbArtists, err
}

func (t *ArtistService) addMissingFeaturingArtists(err error, spotifyIds []string, timeStamp time.Time) (map[string]artistData.Artist, error) {
	if err != nil {
		return make(map[string]artistData.Artist, 0), err
	}

	results := make(map[string]artistData.Artist)
	addedDbArtists, err := t.addMissingArtists(spotifyIds, timeStamp)
	for _, dbArtist := range addedDbArtists {
		if err == nil {
			results[dbArtist.SpotifyId] = dbArtist
		}
	}

	return results, err
}

func (t *ArtistService) findAndAddMissingReleases(err error, currentManager labels.ManagerContext, dbArtist artistData.Artist, timeStamp time.Time) error {
	if err != nil {
		return err
	}

	missingReleases := t.releaseService.GetMissingReleases(dbArtist)

	artistSpotifyIds := t.getFeaturingArtistSpotifyIds(missingReleases)
	existingArtists, err := t.getExistingFeaturingArtists(dbArtist, artistSpotifyIds, timeStamp)

	missingFeaturingArtistsSpotifyIds, err := t.getMissingFeaturingArtistsSpotifyIds(err, existingArtists, artistSpotifyIds)
	addedArtists, err := t.addMissingFeaturingArtists(err, missingFeaturingArtistsSpotifyIds, timeStamp)
	if err != nil {
		return err
	}

	artists := existingArtists
	for key, artist := range addedArtists {
		artists[key] = artist
	}

	return t.releaseService.Add(currentManager, artists, missingReleases, timeStamp)
}

func (t *ArtistService) getDbArtistBySpotifyId(err error, spotifyId string) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	var dbArtist artistData.Artist
	err = data.DB.Model(&artistData.Artist{}).Preload("Releases").Where("spotify_id = ?", spotifyId).FirstOrInit(&dbArtist).Error
	if err != nil {
		t.logger.LogFatal(err, err.Error())
	}

	return dbArtist, err
}

func (t *ArtistService) getExistingFeaturingArtists(dbArtist artistData.Artist, spotifyIds []string, timeStamp time.Time) (map[string]artistData.Artist, error) {
	var existedArtists []artistData.Artist
	err := data.DB.Where("spotify_id IN ?", spotifyIds).Find(&existedArtists).Error
	if err != nil {
		t.logger.LogFatal(err, err.Error())
		return make(map[string]artistData.Artist), err
	}

	results := make(map[string]artistData.Artist, 0)
	for _, artist := range existedArtists {
		results[artist.SpotifyId] = artist
	}

	return results, nil
}

func (t *ArtistService) getInternal(err error, currentManager labels.ManagerContext, id int) (artistModels.Artist, error) {
	if err != nil {
		return artistModels.Artist{}, err
	}

	cacheKey := fmt.Sprintf("%v", id)
	value, isCached := t.cache.TryGet(cacheKey)
	if isCached {
		return value, nil
	}

	var dbArtist artistData.Artist
	err = data.DB.Model(&artistData.Artist{}).Preload("Releases").First(&dbArtist, id).Error
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	err = validators.CurrentDbArtistBelongsToLabel(err, dbArtist, currentManager.LabelId)
	if err != nil {
		return artistModels.Artist{}, err
	}

	artist, err := converters.ToArtist(dbArtist)
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	t.cache.Set(cacheKey, artist, time.Minute*1)

	return artist, err
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

func (t *ArtistService) getFeaturingArtistSpotifyIds(releases []releases.Release) []string {
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

	return spotifyIds
}

func (t *ArtistService) updateLabelIfNeeded(err error, dbArtist artistData.Artist, labelId int) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	if dbArtist.LabelId == 0 {
		dbArtist.LabelId = labelId
		err = data.DB.Save(&dbArtist).Error
	}

	return dbArtist, err
}
