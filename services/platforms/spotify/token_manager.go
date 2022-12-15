package spotify

import (
	"encoding/base64"
	spotifyModels "main/models/platforms/spotify/tokens"
	"net/http"
	"net/url"
	"strings"
	"time"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
)

func getAccessToken(logger logger.Logger, config *spotifyModels.SpotifyClientConfig) (string, error) {
	if len(_tokenContainer.Token) != 0 && time.Now().UTC().Before(_tokenContainer.Expired) {
		return _tokenContainer.Token, nil
	}

	request, err := getAccessTokenRequest(logger, config)
	if err != nil {
		logger.LogError(err, err.Error())
		return "", err
	}

	httpClient := httpClient.New[spotifyModels.SpotifyAccessToken](httpClient.DefaultConfig(logger))
	accessToken, err := httpClient.MakeRequest(request)
	if err != nil {
		logger.LogError(err, err.Error())
		return "", err
	}

	_tokenContainer = spotifyModels.SpotifyAccessTokenContainer{
		Expired: time.Now().Add(time.Second*time.Duration(accessToken.ExpiresIn) - ACCESS_TOKEN_SAFITY_THRESHOLD).UTC(),
		Token:   accessToken.Token,
	}

	return _tokenContainer.Token, nil
}

func getAccessTokenRequest(logger logger.Logger, config *spotifyModels.SpotifyClientConfig) (*http.Request, error) {
	payload := url.Values{}
	payload.Add("grant_type", "client_credentials")

	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(payload.Encode()))
	if err != nil {
		logger.LogError(err, err.Error())
		return nil, err
	}

	credentials := "Basic " + base64.StdEncoding.EncodeToString([]byte(config.ClientId+":"+config.ClientSecret))

	request.Header.Add("Authorization", credentials)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

var _tokenContainer spotifyModels.SpotifyAccessTokenContainer

const ACCESS_TOKEN_SAFITY_THRESHOLD = time.Second * 5
