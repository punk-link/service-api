package repositories

import (
	artistData "main/data/artists"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type PresentationConfigRepositoryService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewPresentationConfigRepository(injector *do.Injector) (PresentationConfigRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &PresentationConfigRepositoryService{
		db:     db,
		logger: logger,
	}, nil
}

func (t *PresentationConfigRepositoryService) Get(err error, artistId int) (artistData.ArtistPresentationConfig, error) {
	if err != nil {
		return artistData.ArtistPresentationConfig{}, err
	}

	var config artistData.ArtistPresentationConfig
	err = t.db.Model(artistData.ArtistPresentationConfig{}).
		First(&config, artistId).
		Error

	return config, t.handleError(err)
}

func (t *PresentationConfigRepositoryService) handleError(err error) error {
	if err != nil {
		t.logger.LogError(err, err.Error())
	}

	return err
}
