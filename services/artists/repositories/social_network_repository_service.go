package repositories

import (
	artistData "main/data/artists"
	"main/helpers"

	"github.com/punk-link/logger"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type SocialNetworkRepositoryService struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewSocialNetworkRepository(injector *do.Injector) (SocialNetworkRepository, error) {
	db := do.MustInvoke[*gorm.DB](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &SocialNetworkRepositoryService{
		db:     db,
		logger: logger,
	}, nil
}

func (t *SocialNetworkRepositoryService) Get(err error, artistId int) []artistData.ArtistSocialNetwork {
	networks := make([]artistData.ArtistSocialNetwork, 0)
	if err != nil {
		return make([]artistData.ArtistSocialNetwork, 0)
	}

	err = t.db.Where("artist_id = ?", artistId).
		Find(&networks).
		Error

	t.handleError(err)
	return networks
}

func (t *SocialNetworkRepositoryService) handleError(err error) error {
	if helpers.ShouldHandleDbError(err) {
		t.logger.LogError(err, err.Error())
	}

	return err
}
