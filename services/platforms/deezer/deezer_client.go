package deezer

import (
	platformServices "main/services/platforms/base"
	"net/http"

	"github.com/punk-link/logger"
)

func makeRequest[T any](logger *logger.Logger, method string, url string, result *T) error {
	request, err := getRequest(method, url)
	if err != nil {
		logger.LogWarn("can't build an http request: %s", err.Error())
		return err
	}

	return platformServices.MakeRequest(logger, request, result)
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
