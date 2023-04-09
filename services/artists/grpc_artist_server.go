package artists

import presentationContracts "github.com/punk-link/presentation-contracts"

type GrpcArtistServer interface {
	GetOne(request *presentationContracts.ArtistRequest) (*presentationContracts.Artist, error)
}
