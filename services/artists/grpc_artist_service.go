package artists

import (
	artistModels "main/models/artists"
	"main/services/artists/converters"
	"main/services/artists/repositories"

	presentationContracts "github.com/punk-link/presentation-contracts"

	"github.com/samber/do"
)

type GrpcArtistService struct {
	artistRepository          repositories.ArtistRepository
	presentationConfigService PresentationConfigServer
	releaseRepository         repositories.ReleaseRepository
}

func NewGrpcArtistService(injector *do.Injector) (GrpcArtistServer, error) {
	artistRepository := do.MustInvoke[repositories.ArtistRepository](injector)
	presentationConfigService := do.MustInvoke[PresentationConfigServer](injector)
	releaseRepository := do.MustInvoke[repositories.ReleaseRepository](injector)

	return &GrpcArtistService{
		artistRepository:          artistRepository,
		presentationConfigService: presentationConfigService,
		releaseRepository:         releaseRepository,
	}, nil
}

func (t *GrpcArtistService) GetOne(request *presentationContracts.ArtistRequest) (*presentationContracts.Artist, error) {
	id := int(request.Id)

	dbArtist, err := t.artistRepository.GetOneSlim(nil, id)
	dbSlimReleases, err := t.releaseRepository.GetSlimByArtistId(err, id)
	presentationConfig, err := t.getPresentationConfig(err, id)
	artist, err := converters.ToArtistMessage(err, dbArtist, presentationConfig)
	slimReleases, err := converters.ToSlimReleaseMessages(err, dbSlimReleases)
	if err != nil {
		return &presentationContracts.Artist{}, err
	}

	artist.Releases = slimReleases
	return artist, nil
}

func (t *GrpcArtistService) getPresentationConfig(err error, artistId int) (artistModels.PresentationConfig, error) {
	if err != nil {
		return artistModels.PresentationConfig{}, err
	}

	return t.presentationConfigService.Get(artistId), nil
}
