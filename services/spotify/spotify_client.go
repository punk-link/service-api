package spotify

import (
	"encoding/json"
	"io"
	"main/services/common"
	"net/http"
)

func MakeRequest[T any](logger *common.Logger, method string, url string, result *T) error {
	request, err := getRequest(logger, method, url)
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.LogError(err, err.Error())
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logger.LogError(err, err.Error())
		return err
	}

	return nil
}

func getRequest(logger *common.Logger, method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "https://api.spotify.com/v1/"+url, nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken(logger)
	if err != nil {
		logger.LogWarn(err.Error())
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
