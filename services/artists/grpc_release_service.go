package artists

import (
	"encoding/json"
	artistModels "main/models/artists"
	"main/services/artists/converters"
	"main/services/artists/repositories"
	platformRepositories "main/services/platforms/repositories"

	"github.com/punk-link/logger"
	presentationContracts "github.com/punk-link/presentation-contracts"

	"github.com/samber/do"
)

type GrpcReleaseService struct {
	artistRepository          repositories.ArtistRepository
	logger                    logger.Logger
	platformUrlRepository     platformRepositories.PlatformUrlRepository
	presentationConfigService PresentationConfigServer
	releaseRepository         repositories.ReleaseRepository
}

func NewGrpcReleaseService(injector *do.Injector) (GrpcReleaseServer, error) {
	artistRepository := do.MustInvoke[repositories.ArtistRepository](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	platformUrlRepository := do.MustInvoke[platformRepositories.PlatformUrlRepository](injector)
	presentationConfigService := do.MustInvoke[PresentationConfigServer](injector)
	releaseRepository := do.MustInvoke[repositories.ReleaseRepository](injector)

	return &GrpcReleaseService{
		artistRepository:          artistRepository,
		logger:                    logger,
		platformUrlRepository:     platformUrlRepository,
		presentationConfigService: presentationConfigService,
		releaseRepository:         releaseRepository,
	}, nil
}

func (t *GrpcReleaseService) GetOne(request *presentationContracts.ReleaseRequest) (*presentationContracts.Release, error) {
	id := int(request.Id)

	dbRelease, err := t.releaseRepository.GetOne(nil, id)
	releaseArtistIds, err := t.unmarshalArtistIds(err, dbRelease.ReleaseArtistIds)
	presentationConfig, err := t.getPresentationConfig(err, releaseArtistIds)
	slimDbArtists, err := t.artistRepository.GetSlim(err, releaseArtistIds)
	platformUrls, err := t.platformUrlRepository.GetByReleaseId(err, id)

	return converters.ToReleaseMessage(err, dbRelease, slimDbArtists, platformUrls, presentationConfig)
}

func (t *GrpcReleaseService) unmarshalArtistIds(err error, idJson string) ([]int, error) {
	if err != nil {
		return make([]int, 0), err
	}

	var results []int
	err = json.Unmarshal([]byte(idJson), &results)
	if err != nil {
		t.logger.LogError(err, err.Error())
		return make([]int, 0), err
	}

	return results, nil
}

func (t *GrpcReleaseService) getPresentationConfig(err error, artistIds []int) (artistModels.PresentationConfig, error) {
	if err != nil || len(artistIds) < 1 {
		return artistModels.PresentationConfig{}, err
	}

	primaryArtistId := artistIds[0]
	return t.presentationConfigService.Get(primaryArtistId), nil
}
