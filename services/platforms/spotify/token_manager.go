package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"main/infrastructure/consul"
	"main/models/platforms/spotify/accessToken"
	"main/services/common"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func getAccessToken(logger *common.Logger) (string, error) {
	if len(tokenContainer.Token) != 0 && time.Now().UTC().Before(tokenContainer.Expired) {
		return tokenContainer.Token, nil
	}

	request, err := getAccessTokenRequest(logger)
	if err != nil {
		logger.LogError(err, err.Error())
		return "", err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.LogError(err, err.Error())
		return "", err
	}
	defer response.Body.Close()

	tokenContainer, err = parseToken(logger, response)
	if err != nil {
		logger.LogError(err, err.Error())
		return "", err
	}

	return tokenContainer.Token, nil
}

func getAccessTokenRequest(logger *common.Logger) (*http.Request, error) {
	payload := url.Values{}
	payload.Add("grant_type", "client_credentials")

	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(payload.Encode()))
	if err != nil {
		logger.LogError(err, err.Error())
		return nil, err
	}

	consul := consul.BuildConsulClient(logger, "service-api")
	spotifySettings := consul.GetOrSet("SpotifySettings", 0).(map[string]interface{})

	clientId := spotifySettings["ClientId"].(string)
	clientSecret := spotifySettings["ClientSecret"].(string)
	credentials := "Basic " + base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret))

	request.Header.Add("Authorization", credentials)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

func parseToken(logger *common.Logger, response *http.Response) (accessToken.SpotifyAccessTokenContainer, error) {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.LogError(err, err.Error())
		return accessToken.SpotifyAccessTokenContainer{}, err
	}

	var newToken accessToken.SpotifyAccessToken
	if err := json.Unmarshal(body, &newToken); err != nil {
		logger.LogError(err, err.Error())
		return accessToken.SpotifyAccessTokenContainer{}, err
	}

	return accessToken.SpotifyAccessTokenContainer{
		Expired: time.Now().Add(time.Second*time.Duration(newToken.ExpiresIn) - safetyThreshold).UTC(),
		Token:   newToken.Token,
	}, nil
}

var tokenContainer accessToken.SpotifyAccessTokenContainer

const safetyThreshold time.Duration = time.Second * 5
