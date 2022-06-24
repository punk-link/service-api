package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"main/models/spotify/accessToken"
	"main/utils"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func getAccessToken() (string, error) {
	if len(tokenContainer.Token) != 0 || time.Now().Before(tokenContainer.ExpiresAt) {
		return tokenContainer.Token, nil
	}

	payload := url.Values{}
	payload.Add("grant_type", "client_credentials")

	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(payload.Encode()))
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	clientId := utils.GetEnvironmentVariable("SPOTIFY_CLIENT_ID")
	clientSecret := utils.GetEnvironmentVariable("SPOTIFY_CLIENT_SECRET")
	credentials := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret))
	request.Header.Add("Authorization", credentials)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", err
	}

	var newToken accessToken.SpotifyAccessToken
	if err := json.Unmarshal(body, &newToken); err != nil {
		log.Error().Err(err).Msg(err.Error())
		return "", err
	}

	tokenContainer = accessToken.SpotifyAccessTokenContainer{
		ExpiresAt: time.Now().Add(time.Second*time.Duration(newToken.ExpiresIn) - safetyThreshold),
		Token:     newToken.Token,
	}

	return tokenContainer.Token, nil
}

var tokenContainer accessToken.SpotifyAccessTokenContainer

const safetyThreshold time.Duration = time.Second * 5
