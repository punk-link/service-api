package spotify

import (
	commonModels "main/models/common"
	"main/services/common"
	platformServices "main/services/platforms/base"
	"net/http"
)

func makeBatchRequestWithSync[T any](logger *common.Logger, method string, syncedUrls []commonModels.SyncedUrl) []commonModels.SyncedResult[T] {
	syncedHttpRequests := make([]commonModels.SyncedHttpRequest, len(syncedUrls))
	for i, syncedUrl := range syncedUrls {
		request, err := getRequest(logger, method, syncedUrl.Url)
		if err != nil {
			logger.LogWarn("can't build an http request: %s", err.Error())
			continue
		}

		syncedHttpRequests[i] = commonModels.SyncedHttpRequest{
			HttpRequest: request,
			Sync:        syncedUrl.Sync,
		}
	}

	return platformServices.MakeBatchRequestWithSync[T](logger, syncedHttpRequests)
}

func makeBatchRequest[T any](logger *common.Logger, method string, urls []string) []T {
	syncedUrls := make([]commonModels.SyncedUrl, len(urls))
	for i, url := range urls {
		syncedUrls[i] = commonModels.SyncedUrl{
			Url: url,
		}
	}

	syncedResults := makeBatchRequestWithSync[T](logger, method, syncedUrls)
	results := make([]T, len(syncedResults))
	for i, result := range syncedResults {
		results[i] = result.Result
	}

	return results
}

func makeRequest[T any](logger *common.Logger, method string, url string, response *T) error {
	request, err := getRequest(logger, method, url)
	if err != nil {
		logger.LogWarn("can't build an http request: %s", err.Error())
		return err
	}

	return platformServices.MakeRequest(logger, request, response)
}

func getRequest(logger *common.Logger, method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "https://api.spotify.com/v1/"+url, nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken(logger)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
