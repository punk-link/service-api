package artists

import presentationContracts "github.com/punk-link/presentation-contracts"

type GrpcReleaseServer interface {
	GetOne(request *presentationContracts.ReleaseRequest) (*presentationContracts.Release, error)
}
