package artists

import (
	artistModels "main/models/artists"
	labelModels "main/models/labels"
)

type SocialNetworkServer interface {
	ArrOrModify(currentManager labelModels.ManagerContext, artistId int, networks []artistModels.SocialNetwork) ([]artistModels.SocialNetwork, error)
	Get(artistId int) []artistModels.SocialNetwork
}
