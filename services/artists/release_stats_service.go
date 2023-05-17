package artists

import (
	artistModels "main/models/artists"
	"main/services/artists/repositories"

	"github.com/samber/do"
)

type ReleaseStatsService struct {
	releaseRepository repositories.ReleaseRepository
}

func NewReleaseStatsService(injector *do.Injector) (ReleaseStatsServer, error) {
	releaseRepository := do.MustInvoke[repositories.ReleaseRepository](injector)

	return &ReleaseStatsService{
		releaseRepository: releaseRepository,
	}, nil
}

func (t *ReleaseStatsService) Get(artistIds []int) artistModels.ArtistReleaseStats {
	if len(artistIds) != 1 {
		return artistModels.ArtistReleaseStats{}
	}

	primaryArtistId := artistIds[0]
	return t.getReleaseStats(primaryArtistId)
}

func (t *ReleaseStatsService) GetOne(artistId int) artistModels.ArtistReleaseStats {
	return artistModels.ArtistReleaseStats{}
}

func (t *ReleaseStatsService) getReleaseStats(artistId int) artistModels.ArtistReleaseStats {
	albumNumber, compilationNumber, singleNumber, _ := t.releaseRepository.GetCountByArtistByType(nil, artistId)

	return artistModels.ArtistReleaseStats{
		AlbumNumber:       albumNumber,
		CompilationNumber: compilationNumber,
		SingleNumber:      singleNumber,
	}
}
