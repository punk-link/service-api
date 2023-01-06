package converters

import (
	presentationContracts "github.com/punk-link/presentation-contracts"
)

func ToMessageFromJson(err error, detailsJson string) (*presentationContracts.ImageDetails, error) {
	if err != nil {
		return &presentationContracts.ImageDetails{}, err
	}

	imageDetails, err := FromJson(detailsJson)
	if err != nil {
		return &presentationContracts.ImageDetails{}, err
	}

	return &presentationContracts.ImageDetails{
		AltText: imageDetails.AltText,
		Height:  int32(imageDetails.Height),
		Url:     imageDetails.Url,
		Width:   int32(imageDetails.Width),
	}, nil
}
