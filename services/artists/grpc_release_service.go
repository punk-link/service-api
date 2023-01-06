package artists

import (
	"main/services/artists/converters"

	presentationContracts "github.com/punk-link/presentation-contracts"

	"github.com/samber/do"
)

type GrpcReleaseService struct {
	artistRepository  *ArtistRepository
	releaseRepository *ReleaseRepository
}

func NewGrpcReleaseService(injector *do.Injector) (*GrpcReleaseService, error) {
	artistRepository := do.MustInvoke[*ArtistRepository](injector)
	releaseRepository := do.MustInvoke[*ReleaseRepository](injector)

	return &GrpcReleaseService{
		artistRepository:  artistRepository,
		releaseRepository: releaseRepository,
	}, nil
}

func (t *GrpcReleaseService) GetOne(request *presentationContracts.ReleaseRequest) (*presentationContracts.Release, error) {
	dbRelease, err := t.releaseRepository.GetOne(nil, int(request.Id))
	releaseArtistIds, err := t.getReleaseArtistIds(err, dbRelease.ReleaseArtistIds)
	slimDbArtists, err := t.artistRepository.GetSlim(err, releaseArtistIds)

	return converters.ToReleaseMessage(err, dbRelease, slimDbArtists)
}

func (t *GrpcReleaseService) getReleaseArtistIds(err error, json string) ([]int, error) {
	if err != nil {
		return make([]int, 0), err
	}

	return make([]int, 0), err
}
