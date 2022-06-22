package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"main/models/spotify/accessToken"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getAccessToken() (string, error) {
	const baseUrl string = "https://accounts.spotify.com/api/token"

	if len(tokenContainer.Token) != 0 || time.Now().Before(tokenContainer.ExpiresAt) {
		return tokenContainer.Token, nil
	}

	payload := url.Values{}
	payload.Add("grant_type", "client_credentials")

	request, err := http.NewRequest("POST", baseUrl, strings.NewReader(payload.Encode()))
	if err != nil {
		return "", err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	credentials := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret))
	request.Header.Add("Authorization", credentials)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var newToken accessToken.SpotifyAccessToken
	if err := json.Unmarshal(body, &newToken); err != nil {
		return "", err
	}

	tokenContainer = accessToken.SpotifyAccessTokenContainer{
		ExpiresAt: time.Now().Add(time.Second*time.Duration(newToken.ExpiresIn) - safetyThreshold),
		Token:     newToken.Token,
	}

	return tokenContainer.Token, nil
}

var tokenContainer accessToken.SpotifyAccessTokenContainer

var clientId = ""
var clientSecret = ""

const safetyThreshold time.Duration = time.Second * 5
