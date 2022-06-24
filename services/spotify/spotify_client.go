package spotify

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

func MakeRequest[T any](method string, url string, result *T) error {
	request, err := http.NewRequest(method, "https://api.spotify.com/v1/"+url, nil)
	if err != nil {
		return err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		log.Warn().Msg(err.Error())
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
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
