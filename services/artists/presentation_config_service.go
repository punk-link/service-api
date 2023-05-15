package artists

import (
	"encoding/json"
	artistData "main/data/artists"
	artistModels "main/models/artists"

	"github.com/punk-link/logger"
	presentationContractsConstants "github.com/punk-link/presentation-contracts/constants"
	"github.com/samber/do"
)

type PresentationConfigService struct {
	logger     logger.Logger
	repository PresentationConfigRepository
}

func NewPresentationConfigService(injector *do.Injector) (PresentationConfigServer, error) {
	logger := do.MustInvoke[logger.Logger](injector)
	repository := do.MustInvoke[PresentationConfigRepository](injector)

	return &PresentationConfigService{
		logger:     logger,
		repository: repository,
	}, nil
}

func (t *PresentationConfigService) Get(artistId int) artistModels.PresentationConfig {
	shareableSocialNetworkIds := t.getShareableSocialNetworkIds(artistId)

	return artistModels.PresentationConfig{
		ShareableSocialNetworkIds: shareableSocialNetworkIds,
	}
}

func (t *PresentationConfigService) getShareableSocialNetworkIds(artistId int) []string {
	dbConfig, err := t.repository.Get(nil, artistId)
	if err != nil || (dbConfig == artistData.ArtistPresentationConfig{}) {
		return getDefaultShareableSocialNetworkIds()
	}

	var results []string
	err = json.Unmarshal([]byte(dbConfig.Value), &results)
	if err != nil {
		t.logger.LogError(err, err.Error())
		return getDefaultShareableSocialNetworkIds()
	}

	return results
}

func getDefaultShareableSocialNetworkIds() []string {
	return []string{
		presentationContractsConstants.FACEBOOK,
		// presentationContractsConstants.INSTAGRAM,
		presentationContractsConstants.TELEGRAM,
		presentationContractsConstants.TWITTER,
		presentationContractsConstants.VK,
	}
}
