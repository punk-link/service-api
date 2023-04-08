package converters

import (
	"encoding/json"
	"fmt"
	"main/helpers"
	commonModels "main/models/common"
	spotifyPlatformModels "main/models/platforms/spotify"
)

func FromJson(detailsJson string) (commonModels.ImageDetails, error) {
	imageDetails := commonModels.ImageDetails{}
	var err error

	if detailsJson != emptyJsonToken {
		err = json.Unmarshal([]byte(detailsJson), &imageDetails)
	}

	return imageDetails, err
}

func ToImageDetailsFromSpotify(imageDetails []spotifyPlatformModels.ImageDetails, altText string) commonModels.ImageDetails {
	if len(imageDetails) == 0 {
		return commonModels.ImageDetails{}
	}

	helpers.OrderImageDetailsDesc(imageDetails)
	return commonModels.ImageDetails{
		AltText: altText,
		Height:  imageDetails[0].Height,
		Url:     imageDetails[0].Url,
		Width:   imageDetails[0].Width,
	}
}

func ToJsonFromSpotify(imageDetails []spotifyPlatformModels.ImageDetails, altText string) (string, error) {
	imageDetailsJson := emptyJsonToken
	var err error

	if len(imageDetails) != 0 {
		imageDetails := ToImageDetailsFromSpotify(imageDetails, altText)
		imageDetailsJson, err = ToJson(imageDetails)
	}

	return imageDetailsJson, err
}

func ToJson(imageDetails commonModels.ImageDetails) (string, error) {
	imageDetailsBytes, err := json.Marshal(imageDetails)
	if err != nil {
		return "", fmt.Errorf("can't serialize image details: '%s'", err.Error())
	}

	return string(imageDetailsBytes), nil
}

const emptyJsonToken string = "{}"
