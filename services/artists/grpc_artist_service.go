package artists

import (
	"main/services/artists/converters"

	presentationContracts "github.com/punk-link/presentation-contracts"

	"github.com/samber/do"
)

type GrpcArtistService struct {
	artistRepository  ArtistRepository
	releaseRepository ReleaseRepository
}

func NewGrpcArtistService(injector *do.Injector) (GrpcArtistServer, error) {
	artistRepository := do.MustInvoke[ArtistRepository](injector)
	releaseRepository := do.MustInvoke[ReleaseRepository](injector)

	return &GrpcArtistService{
		artistRepository:  artistRepository,
		releaseRepository: releaseRepository,
	}, nil
}

func (t *GrpcArtistService) GetOne(request *presentationContracts.ArtistRequest) (*presentationContracts.Artist, error) {
	id := int(request.Id)

	dbArtist, err := t.artistRepository.GetOneSlim(nil, id)
	dbSlimReleases, err := t.releaseRepository.GetSlimByArtistId(err, id)
	artist, err := converters.ToArtistMessage(err, dbArtist)
	slimReleases, err := converters.ToSlimReleaseMessages(err, dbSlimReleases)
	if err != nil {
		return &presentationContracts.Artist{}, err
	}

	artist.Releases = slimReleases
	return artist, nil
}
