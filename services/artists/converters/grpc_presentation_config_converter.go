package converters

import (
	artistModels "main/models/artists"

	presentationContracts "github.com/punk-link/presentation-contracts"
)

func ToPresentationConfigMessage(err error, config artistModels.PresentationConfig) (*presentationContracts.PresentationConfig, error) {
	if err != nil {
		return &presentationContracts.PresentationConfig{}, err
	}

	return &presentationContracts.PresentationConfig{
		ShareableSocialNetworkIds: config.ShareableSocialNetworkIds,
	}, nil
}
