package artists

import (
	"errors"
	"main/data"
	artistData "main/data/artists"
	artistModels "main/models/artists"
	"main/models/labels"
	spotifyArtists "main/models/spotify/artists"
	"main/models/spotify/releases"
	"main/services/common"
	"main/services/spotify"
	"time"
)

type ArtistService struct {
	logger         *common.Logger
	spotifyService *spotify.SpotifyService
	releaseService *ReleaseService
}

func BuildArtistService(logger *common.Logger, releaseService *ReleaseService, spotifyService *spotify.SpotifyService) *ArtistService {
	return &ArtistService{
		logger:         logger,
		releaseService: releaseService,
		spotifyService: spotifyService,
	}
}

func (t *ArtistService) AddArtist(currentManager labels.ManagerContext, spotifyId string) (interface{}, error) {
	if spotifyId == "" {
		return "", errors.New("artist's spotify ID is empty")
	}

	dbArtist, err := t.getDatabaseArtistBySpotifyId(spotifyId)
	if err != nil {
		return "", err
	}

	now := time.Now().UTC()
	if dbArtist.Id != 0 {
		if dbArtist.LabelId != currentManager.LabelId {
			return "", errors.New("artist already added to another label")
		}
	} else {
		spotifyArtist := t.spotifyService.GetArtist(spotifyId)
		dbArtist, err = t.addArtist(spotifyArtist, currentManager.LabelId, now)
		if err != nil {
			return "", err
		}
	}

	missingReleaseSpotifyIds := t.releaseService.GetMissingReleaseIds(dbArtist)
	releases := t.spotifyService.GetReleasesDetails(missingReleaseSpotifyIds)
	artists, err := t.GetOrAddReleasesArtists(dbArtist, releases, now)
	if err != nil {
		return make([]artistModels.Release, 0), err
	}

	_, err = t.releaseService.AddReleases(currentManager, artists, releases, now)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (t *ArtistService) GetOrAddReleasesArtists(dbArtist artistData.Artist, releases []releases.Release, timeStamp time.Time) (map[string]artistData.Artist, error) {
	spotifyIds := t.getReleasesArtistsSpotifyIds(releases)
	results, err := t.getDatabaseArtistsBySpotifyIds(spotifyIds)
	if err != nil {
		return make(map[string]artistData.Artist, 0), err
	}

	existingSpotifyIds := make(map[string]int, len(results))
	for _, artist := range results {
		existingSpotifyIds[artist.SpotifyId] = 0
	}

	dbArtists, err := t.addMissingArtists(existingSpotifyIds, spotifyIds, timeStamp)
	for _, dbArtist := range dbArtists {
		if err == nil {
			results[dbArtist.SpotifyId] = dbArtist
		}
	}

	return results, nil
}

func (t *ArtistService) SearchArtist(query string) []artistModels.ArtistSearchResult {
	var result []artistModels.ArtistSearchResult

	const minimalQueryLength int = 3
	if len(query) < minimalQueryLength {
		return result
	}

	results := t.spotifyService.SearchArtist(query)

	return spotify.ToArtistSearchResults(results)
}

func (t *ArtistService) addMissingArtists(existingSpotifyIds map[string]int, spotifyIds []string, timeStamp time.Time) ([]artistData.Artist, error) {
	missingSpotifyIds := make([]string, 0)
	for _, id := range spotifyIds {
		if _, isExists := existingSpotifyIds[id]; !isExists {
			missingSpotifyIds = append(missingSpotifyIds, id)
		}
	}

	missingArtists := t.spotifyService.GetArtists(missingSpotifyIds)

	dbArtists := make([]artistData.Artist, len(missingArtists))
	for i, artist := range missingArtists {
		dbArtists[i] = t.toDbArtist(artist, 0, timeStamp)
	}

	result := data.DB.CreateInBatches(&dbArtists, 50)
	if result.Error != nil {
		t.logger.LogError(result.Error, result.Error.Error())
		return make([]artistData.Artist, 0), result.Error
	}

	return dbArtists, nil
}

func (t *ArtistService) addArtist(artist spotifyArtists.Artist, labelId int, timeStamp time.Time) (artistData.Artist, error) {
	dbArtist := t.toDbArtist(artist, labelId, timeStamp)

	result := data.DB.Create(&dbArtist)
	if result.Error != nil {
		t.logger.LogError(result.Error, result.Error.Error())
		return artistData.Artist{}, result.Error
	}

	return dbArtist, nil
}

func (t *ArtistService) getDatabaseArtistsBySpotifyIds(spotifyIds []string) (map[string]artistData.Artist, error) {
	var existedArtists []artistData.Artist
	queryResult := data.DB.Where("spotify_id IN ?", spotifyIds).Find(&existedArtists)
	if queryResult.Error != nil {
		t.logger.LogFatal(queryResult.Error, queryResult.Error.Error())
		return make(map[string]artistData.Artist), queryResult.Error
	}

	results := make(map[string]artistData.Artist, 0)
	for _, artist := range existedArtists {
		results[artist.SpotifyId] = artist
	}

	return results, nil
}

func (t *ArtistService) getDatabaseArtistBySpotifyId(spotifyId string) (artistData.Artist, error) {
	var dbArtist artistData.Artist
	queryResult := data.DB.Model(&artistData.Artist{}).Preload("Releases").Where("spotify_id = ?", spotifyId).FirstOrInit(&dbArtist)
	if queryResult.Error != nil {
		t.logger.LogFatal(queryResult.Error, queryResult.Error.Error())
		return dbArtist, queryResult.Error
	}

	return dbArtist, nil
}

func (t *ArtistService) getReleasesArtistsSpotifyIds(releases []releases.Release) []string {
	artistIds := make(map[string]int)
	for _, release := range releases {
		for _, artist := range release.Artists {
			if _, isExists := artistIds[artist.Id]; !isExists {
				artistIds[artist.Id] = 0
			}
		}
	}

	spotifyIds := make([]string, 0)
	for i := range artistIds {
		spotifyIds = append(spotifyIds, i)
	}

	return spotifyIds
}

func (t *ArtistService) toDbArtist(artist spotifyArtists.Artist, labelId int, timeStamp time.Time) artistData.Artist {
	return artistData.Artist{
		Created:   timeStamp,
		LabelId:   labelId,
		Name:      artist.Name,
		SpotifyId: artist.Id,
		Updated:   timeStamp,
	}
}
