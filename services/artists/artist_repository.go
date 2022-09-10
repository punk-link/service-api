package artists

import (
	"main/data"
	artistData "main/data/artists"
	"main/services/common"
)

func createDbArtist(logger *common.Logger, err error, artist *artistData.Artist) error {
	if err != nil {
		return err
	}

	err = data.DB.Create(artist).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func createDbArtistsInBatches(logger *common.Logger, err error, artists *[]artistData.Artist) error {
	if err != nil || len(*artists) == 0 {
		return err
	}

	err = data.DB.CreateInBatches(&artists, CREATE_ARTISTS_BATCH_SIZE).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func getDbArtist(logger *common.Logger, err error, id int) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	var artist artistData.Artist
	err = data.DB.Model(&artistData.Artist{}).
		First(&artist, id).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artist, err
}

func getDbArtists(logger *common.Logger, err error, ids []int) ([]artistData.Artist, error) {
	if err != nil {
		return make([]artistData.Artist, 0), err
	}

	var artists []artistData.Artist
	err = data.DB.Model(&artistData.Artist{}).
		Find(&artists, ids).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artists, err
}

func getDbArtistBySpotifyId(logger *common.Logger, err error, spotifyId string) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	var artist artistData.Artist
	err = data.DB.Model(&artistData.Artist{}).
		Where("spotify_id = ?", spotifyId).
		FirstOrInit(&artist).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artist, err
}

func getDbArtistsBySpotifyIds(logger *common.Logger, err error, spotifyIds []string) ([]artistData.Artist, error) {
	if err != nil {
		return make([]artistData.Artist, 0), err
	}

	var artists []artistData.Artist
	err = data.DB.Where("spotify_id IN ?", spotifyIds).
		Find(&artists).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artists, err
}

func updateDbArtist(logger *common.Logger, err error, artist *artistData.Artist) error {
	if err != nil {
		return err
	}

	err = data.DB.Save(&artist).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_ARTISTS_BATCH_SIZE int = 50
