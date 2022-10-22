package spotify

import (
	commonModels "main/models/common"
	"main/models/platforms/spotify/accessToken"
	"net/http"

	httpClient "github.com/punk-link/http-client"
	"github.com/punk-link/logger"
)

func makeBatchRequestWithSync[T any](logger logger.Logger, config *accessToken.SpotifyClientConfig, method string, syncedUrls []commonModels.SyncedUrl) []httpClient.SyncedResult[T] {
	syncedHttpRequests := make([]httpClient.SyncedRequest, len(syncedUrls))
	for i, syncedUrl := range syncedUrls {
		request, err := getRequest(logger, config, method, syncedUrl.Url)
		if err != nil {
			logger.LogWarn("can't build an http request: %s", err.Error())
			continue
		}

		syncedHttpRequests[i] = httpClient.SyncedRequest{
			HttpRequest: request,
			SyncKey:     syncedUrl.Sync,
		}
	}

	return httpClient.MakeBatchRequestWithSyncKeys[T](httpClient.DefaultConfig(logger), syncedHttpRequests)
}

func makeBatchRequest[T any](logger logger.Logger, config *accessToken.SpotifyClientConfig, method string, urls []string) []T {
	syncedUrls := make([]commonModels.SyncedUrl, len(urls))
	for i, url := range urls {
		syncedUrls[i] = commonModels.SyncedUrl{
			Url: url,
		}
	}

	syncedResults := makeBatchRequestWithSync[T](logger, config, method, syncedUrls)
	results := make([]T, len(syncedResults))
	for i, result := range syncedResults {
		results[i] = result.Result
	}

	return results
}

func makeRequest[T any](logger logger.Logger, config *accessToken.SpotifyClientConfig, method string, url string, response *T) error {
	request, err := getRequest(logger, config, method, url)
	if err != nil {
		logger.LogWarn("can't build an http request: %s", err.Error())
		return err
	}

	return httpClient.MakeRequest(httpClient.DefaultConfig(logger), request, response)
}

func getRequest(logger logger.Logger, config *accessToken.SpotifyClientConfig, method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "https://api.spotify.com/v1/"+url, nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken(logger, config)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
