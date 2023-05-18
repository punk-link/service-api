package artists

import artistModels "main/models/artists"

type SocialNetworkServer interface {
	Get(artistId int) []artistModels.SocialNetwork
}
