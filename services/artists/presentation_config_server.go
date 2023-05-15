package artists

import artistModels "main/models/artists"

type PresentationConfigServer interface {
	Get(artistId int) artistModels.PresentationConfig
}
