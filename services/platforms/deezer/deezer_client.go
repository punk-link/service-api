package deezer

import (
	commonModels "main/models/common"
	"main/services/common"
	platformServices "main/services/platforms/base"
	"net/http"
)

func makeBatchRequestWithSync[T any](logger *common.Logger, method string, syncedUrls []commonModels.SyncedUrl) []commonModels.SyncedResult[T] {
	syncedHttpRequests := make([]commonModels.SyncedHttpRequest, len(syncedUrls))
	for i, syncedUrl := range syncedUrls {
		request, err := getRequest(method, syncedUrl.Url)
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

func getRequest(method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "http://api.deezer.com/"+url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}
