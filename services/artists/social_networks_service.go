package artists

import (
	artistModels "main/models/artists"
	labelModels "main/models/labels"
	repositories "main/services/artists/repositories"

	"github.com/samber/do"
)

type SocialNetworksService struct {
	repository repositories.SocialNetworkRepository
}

func NewSocialNetworksService(injector *do.Injector) (SocialNetworkServer, error) {
	repository := do.MustInvoke[repositories.SocialNetworkRepository](injector)

	return &SocialNetworksService{
		repository: repository,
	}, nil
}

func (t *SocialNetworksService) ArrOrModify(currentManager labelModels.ManagerContext, artistId int, networks []artistModels.SocialNetwork) ([]artistModels.SocialNetwork, error) {
	// TODO: check access

	// TODO: add

	return t.Get(artistId), nil
}

func (t *SocialNetworksService) Get(artistId int) []artistModels.SocialNetwork {
	dbNetworks := t.repository.Get(nil, artistId)

	results := make([]artistModels.SocialNetwork, len(dbNetworks))
	for i, network := range dbNetworks {
		results[i] = artistModels.SocialNetwork{
			Id:  network.NetworkId,
			Url: network.Url,
		}
	}

	return results
}
