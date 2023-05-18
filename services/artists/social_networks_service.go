package artists

import (
	artistModels "main/models/artists"

	"github.com/samber/do"
)

type SocialNetworksService struct{}

func NewSocialNetworksService(injector *do.Injector) (SocialNetworkServer, error) {
	return &SocialNetworksService{}, nil
}

func (t *SocialNetworksService) Get(artistId int) []artistModels.SocialNetwork {
	return make([]artistModels.SocialNetwork, 0)
}
