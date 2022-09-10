package converters

import (
	"encoding/json"
	"fmt"
	"main/helpers"
	"main/models/common"
	"main/models/spotify"
)

func FromJson(detailsJson string) (common.ImageDetails, error) {
	imageDetails := common.ImageDetails{}
	var err error

	if detailsJson != emptyJsonToken {
		err = json.Unmarshal([]byte(detailsJson), &imageDetails)
	}

	return imageDetails, err
}

func ToImageDetailsFromSpotify(imageDetails []spotify.ImageDetails, altText string) common.ImageDetails {
	if len(imageDetails) == 0 {
		return common.ImageDetails{}
	}

	helpers.OrderImageDetailsDesc(imageDetails)
	return common.ImageDetails{
		AltText: altText,
		Height:  imageDetails[0].Height,
		Url:     imageDetails[0].Url,
		Width:   imageDetails[0].Width,
	}
}

func ToJsonFromSpotify(imageDetails []spotify.ImageDetails, altText string) (string, error) {
	imageDetailsJson := emptyJsonToken
	var err error

	if len(imageDetails) != 0 {
		imageDetails := ToImageDetailsFromSpotify(imageDetails, altText)
		imageDetailsJson, err = ToJson(imageDetails)
	}

	return imageDetailsJson, err
}

func ToJson(imageDetails common.ImageDetails) (string, error) {
	imageDetailsBytes, err := json.Marshal(imageDetails)
	if err != nil {
		return "", fmt.Errorf("can't serialize image details: '%s'", err.Error())
	}

	return string(imageDetailsBytes), nil
}

const emptyJsonToken string = "{}"
