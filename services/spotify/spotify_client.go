package spotify

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func MakeRequest[T any](method string, url string, result *T) error {
	const baseUrl string = "https://api.spotify.com/v1/"

	request, err := http.NewRequest(method, baseUrl+url, nil)
	if err != nil {
		return err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	return nil
}
