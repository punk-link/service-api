package repositories

import (
	artistData "main/data/artists"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type ArtistRepositoryService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewArtistRepository(injector *do.Injector) (ArtistRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &ArtistRepositoryService{
		db:     db,
		logger: logger,
	}, nil
}

func (t *ArtistRepositoryService) Create(err error, artist *artistData.Artist) error {
	if err != nil {
		return err
	}

	err = t.db.Create(artist).Error
	return t.handleError(err)
}

func (t *ArtistRepositoryService) CreateInBatches(err error, artists *[]artistData.Artist) error {
	if err != nil || len(*artists) == 0 {
		return err
	}

	err = t.db.CreateInBatches(&artists, CREATE_ARTISTS_BATCH_SIZE).Error
	return t.handleError(err)
}

func (t *ArtistRepositoryService) Get(err error, ids []int) ([]artistData.Artist, error) {
	if err != nil {
		return make([]artistData.Artist, 0), err
	}

	if len(ids) == 0 {
		return make([]artistData.Artist, 0), err
	}

	var artists []artistData.Artist
	err = t.db.Model(&artistData.Artist{}).
		Find(&artists, ids).
		Error

	return artists, t.handleError(err)
}

func (t *ArtistRepositoryService) GetSlim(err error, ids []int) ([]artistData.SlimArtist, error) {
	if err != nil {
		return make([]artistData.SlimArtist, 0), err
	}

	if len(ids) == 0 {
		return make([]artistData.SlimArtist, 0), err
	}

	var artists []artistData.SlimArtist
	err = t.db.Model(&artistData.Artist{}).
		Find(&artists, ids).
		Error

	return artists, t.handleError(err)
}

func (t *ArtistRepositoryService) GetBySpotifyIds(err error, spotifyIds []string) ([]artistData.Artist, error) {
	if err != nil {
		return make([]artistData.Artist, 0), err
	}

	var artists []artistData.Artist
	err = t.db.Where("spotify_id IN ?", spotifyIds).
		Find(&artists).
		Error

	return artists, t.handleError(err)
}

func (t *ArtistRepositoryService) GetIdsByLabelId(err error, labelId int) ([]int, error) {
	if err != nil {
		return make([]int, 0), err
	}

	var artistIds []int
	err = t.db.Model(&artistData.Artist{}).
		Select("id").
		Where("label_id = ?", labelId).
		Find(&artistIds).
		Error

	return artistIds, t.handleError(err)
}

func (t *ArtistRepositoryService) GetOne(err error, id int) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	var artist artistData.Artist
	err = t.db.Model(&artistData.Artist{}).
		First(&artist, id).
		Error

	return artist, t.handleError(err)
}

func (t *ArtistRepositoryService) GetOneBySpotifyId(err error, spotifyId string) (artistData.Artist, error) {
	if err != nil {
		return artistData.Artist{}, err
	}

	var artist artistData.Artist
	err = t.db.Model(&artistData.Artist{}).
		Where("spotify_id = ?", spotifyId).
		FirstOrInit(&artist).
		Error

	return artist, t.handleError(err)
}

func (t *ArtistRepositoryService) GetOneSlim(err error, id int) (artistData.SlimArtist, error) {
	if err != nil {
		return artistData.SlimArtist{}, err
	}

	var artist artistData.SlimArtist
	err = t.db.Model(&artistData.Artist{}).
		First(&artist, id).
		Error

	return artist, t.handleError(err)
}

func (t *ArtistRepositoryService) Update(err error, artist *artistData.Artist) error {
	if err != nil {
		return err
	}

	err = t.db.Save(&artist).Error
	return t.handleError(err)
}

func (t *ArtistRepositoryService) handleError(err error) error {
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}

const CREATE_ARTISTS_BATCH_SIZE = 50
