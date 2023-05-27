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
	releaseStatsService       ReleaseStatsServer
	socialNetworksService     SocialNetworkServer
}

func NewGrpcArtistService(injector *do.Injector) (GrpcArtistServer, error) {
	artistRepository := do.MustInvoke[repositories.ArtistRepository](injector)
	presentationConfigService := do.MustInvoke[PresentationConfigServer](injector)
	releaseRepository := do.MustInvoke[repositories.ReleaseRepository](injector)
	releaseStatsService := do.MustInvoke[ReleaseStatsServer](injector)
	socialNetworksService := do.MustInvoke[SocialNetworkServer](injector)

	return &GrpcArtistService{
		artistRepository:          artistRepository,
		presentationConfigService: presentationConfigService,
		releaseRepository:         releaseRepository,
		releaseStatsService:       releaseStatsService,
		socialNetworksService:     socialNetworksService,
	}, nil
}

func (t *GrpcArtistService) GetOne(request *presentationContracts.ArtistRequest) (*presentationContracts.Artist, error) {
	id := int(request.Id)

	slimReleasesChan := make(chan []*presentationContracts.SlimRelease, 1)
	go func() {
		dbSlimReleases, err := t.releaseRepository.GetSlimByArtistId(nil, id)
		slimReleases, err := converters.ToSlimReleaseMessages(err, dbSlimReleases)
		if err == nil {
			slimReleasesChan <- slimReleases
		}

		close(slimReleasesChan)
	}()

	dbArtist, err := t.artistRepository.GetOneSlim(nil, id)
	presentationConfig, err := t.getPresentationConfig(err, id)
	releaseStats, err := t.getReleaseStats(err, id)
	socialNetworks, err := t.getSocialNetworks(err, id)

	artist, err := converters.ToArtistMessage(err, dbArtist, releaseStats, socialNetworks, presentationConfig)
	if err != nil {
		return &presentationContracts.Artist{}, err
	}

	artist.Releases = <-slimReleasesChan
	return artist, nil
}

func (t *GrpcArtistService) getPresentationConfig(err error, artistId int) (artistModels.PresentationConfig, error) {
	if err != nil {
		return artistModels.PresentationConfig{}, err
	}

	return t.presentationConfigService.Get(artistId), nil
}

func (t *GrpcArtistService) getReleaseStats(err error, artistId int) (artistModels.ArtistReleaseStats, error) {
	if err != nil {
		return artistModels.ArtistReleaseStats{}, err
	}

	return t.releaseStatsService.GetOne(artistId), nil
}

func (t *GrpcArtistService) getSocialNetworks(err error, artistId int) ([]artistModels.SocialNetwork, error) {
	if err != nil {
		return make([]artistModels.SocialNetwork, 0), err
	}

	return t.socialNetworksService.Get(artistId), nil
}
