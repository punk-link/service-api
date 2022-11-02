package artists

import (
	artistData "main/data/artists"

	"github.com/punk-link/logger"
	"gorm.io/gorm"
)

func createDbArtist(db *gorm.DB, logger logger.Logger, err error, artist *artistData.Artist) error {
	if err != nil {
		return err
	}

	err = db.Create(artist).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func createDbArtistsInBatches(db *gorm.DB, logger logger.Logger, err error, artists *[]artistData.Artist) error {
	if err != nil || len(*artists) == 0 {
		return err
	}

	err = db.CreateInBatches(&artists, CREATE_ARTISTS_BATCH_SIZE).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

func getDbArtist(db *gorm.DB, logger logger.Logger, err error, id int) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	var artist artistData.Artist
	err = db.Model(&artistData.Artist{}).
		First(&artist, id).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artist, err
}

func getDbArtistIdsByLabelId(db *gorm.DB, logger logger.Logger, err error, labelId int) ([]int, error) {
	if err != nil {
		return make([]int, 0), err
	}

	var artistIds []int
	err = db.Model(&artistData.Artist{}).
		Select("id").
		Where("label_id = ?", labelId).
		Find(&artistIds).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artistIds, err
}

func getDbArtists(db *gorm.DB, logger logger.Logger, err error, ids []int) ([]artistData.Artist, error) {
	if err != nil {
		return make([]artistData.Artist, 0), err
	}

	var artists []artistData.Artist
	err = db.Model(&artistData.Artist{}).
		Find(&artists, ids).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artists, err
}

func getDbArtistBySpotifyId(db *gorm.DB, logger logger.Logger, err error, spotifyId string) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	var artist artistData.Artist
	err = db.Model(&artistData.Artist{}).
		Where("spotify_id = ?", spotifyId).
		FirstOrInit(&artist).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artist, err
}

func getDbArtistsBySpotifyIds(db *gorm.DB, logger logger.Logger, err error, spotifyIds []string) ([]artistData.Artist, error) {
	if err != nil {
		return make([]artistData.Artist, 0), err
	}

	var artists []artistData.Artist
	err = db.Where("spotify_id IN ?", spotifyIds).
		Find(&artists).
		Error

	if err != nil {
		logger.LogError(err, err.Error())
	}

	return artists, err
}

func updateDbArtist(db *gorm.DB, logger logger.Logger, err error, artist *artistData.Artist) error {
	if err != nil {
		return err
	}

	err = db.Save(&artist).Error
	if err != nil {
		logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_ARTISTS_BATCH_SIZE = 50
