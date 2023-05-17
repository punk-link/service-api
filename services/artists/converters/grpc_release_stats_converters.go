package converters

import (
	artistModels "main/models/artists"

	presentationContracts "github.com/punk-link/presentation-contracts"
)

func ToReleaseStatsMessage(err error, releaseStats artistModels.ArtistReleaseStats) (*presentationContracts.ReleaseStats, error) {
	if err != nil {
		return &presentationContracts.ReleaseStats{}, err
	}

	return &presentationContracts.ReleaseStats{
		AlbumNumber:       int32(releaseStats.AlbumNumber),
		CompilationNumber: int32(releaseStats.CompilationNumber),
		SingleNumber:      int32(releaseStats.SingleNumber),
	}, nil
}
