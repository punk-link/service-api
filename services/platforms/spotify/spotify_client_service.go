package spotify

import (
	"encoding/base64"
	"fmt"
	tokenSpotifyPlatformModels "main/models/platforms/spotify/tokens"
	"net/http"
	"net/url"
	"strings"
	"time"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

type SpotifyClientService struct {
	config     *tokenSpotifyPlatformModels.SpotifyClientConfig
	httpConfig *httpClient.HttpClientConfig
	logger     logger.Logger
}

func NewSpotifyClient(injector *do.Injector) (SpotifyClient, error) {
	config := do.MustInvoke[*tokenSpotifyPlatformModels.SpotifyClientConfig](injector)
	httpConfig := do.MustInvoke[*httpClient.HttpClientConfig](injector)
	logger := do.MustInvoke[logger.Logger](injector)

	return &SpotifyClientService{
		config:     config,
		httpConfig: httpConfig,
		logger:     logger,
	}, nil
}

func (t *SpotifyClientService) Request(params [][]string, method string, url string) []*http.Request {
	httpRequests := make([]*http.Request, len(params))
	for i, param := range params {
		joinedParams := strings.Join(param, ",")
		request, err := t.RequestOne(method, fmt.Sprintf(url, joinedParams))
		if err != nil {
			t.logger.LogWarn("can't build a http request: %s", err.Error())
			continue
		}

		httpRequests[i] = request
	}

	return httpRequests
}

func (t *SpotifyClientService) RequestOne(method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "https://api.spotify.com/v1/"+url, nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := t.getAccessToken()
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func (t *SpotifyClientService) getAccessToken() (string, error) {
	if len(_tokenContainer.Token) != 0 && time.Now().UTC().Before(_tokenContainer.Expired) {
		return _tokenContainer.Token, nil
	}

	request, err := t.getAccessTokenRequest()
	if err != nil {
		t.logger.LogError(err, err.Error())
		return "", err
	}

	httpClient := httpClient.New[tokenSpotifyPlatformModels.SpotifyAccessToken](t.httpConfig)
	accessToken, err := httpClient.MakeRequest(request)
	if err != nil {
		t.logger.LogError(err, err.Error())
		return "", err
	}

	_tokenContainer = tokenSpotifyPlatformModels.SpotifyAccessTokenContainer{
		Expired: time.Now().Add(time.Second*time.Duration(accessToken.ExpiresIn) - ACCESS_TOKEN_SAFITY_THRESHOLD).UTC(),
		Token:   accessToken.Token,
	}

	return _tokenContainer.Token, nil
}

func (t *SpotifyClientService) getAccessTokenRequest() (*http.Request, error) {
	payload := url.Values{}
	payload.Add("grant_type", "client_credentials")

	request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(payload.Encode()))
	if err != nil {
		t.logger.LogError(err, err.Error())
		return nil, err
	}

	credentials := "Basic " + base64.StdEncoding.EncodeToString([]byte(t.config.ClientId+":"+t.config.ClientSecret))

	request.Header.Add("Authorization", credentials)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request, nil
}

var _tokenContainer tokenSpotifyPlatformModels.SpotifyAccessTokenContainer

const ACCESS_TOKEN_SAFITY_THRESHOLD = time.Second * 5
