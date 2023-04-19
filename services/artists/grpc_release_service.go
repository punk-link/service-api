package artists

import (
	"encoding/json"
	"main/services/artists/converters"
	platformRepositories "main/services/platforms/repositories"

	"github.com/punk-link/logger"
	presentationContracts "github.com/punk-link/presentation-contracts"

	"github.com/samber/do"
)

type GrpcReleaseService struct {
	artistRepository      ArtistRepository
	logger                logger.Logger
	platformUrlRepository platformRepositories.PlatformUrlRepository
	releaseRepository     ReleaseRepository
}

func NewGrpcReleaseService(injector *do.Injector) (GrpcReleaseServer, error) {
	artistRepository := do.MustInvoke[ArtistRepository](injector)
	logger := do.MustInvoke[logger.Logger](injector)
	platformUrlRepository := do.MustInvoke[platformRepositories.PlatformUrlRepository](injector)
	releaseRepository := do.MustInvoke[ReleaseRepository](injector)

	return &GrpcReleaseService{
		artistRepository:      artistRepository,
		logger:                logger,
		platformUrlRepository: platformUrlRepository,
		releaseRepository:     releaseRepository,
	}, nil
}

func (t *GrpcReleaseService) GetOne(request *presentationContracts.ReleaseRequest) (*presentationContracts.Release, error) {
	id := int(request.Id)

	dbRelease, err := t.releaseRepository.GetOne(nil, id)
	releaseArtistIds, err := t.unmarshalArtistIds(err, dbRelease.ReleaseArtistIds)
	slimDbArtists, err := t.artistRepository.GetSlim(err, releaseArtistIds)
	platformUrls, err := t.platformUrlRepository.GetByReleaseId(err, id)

	return converters.ToReleaseMessage(err, dbRelease, slimDbArtists, platformUrls)
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
