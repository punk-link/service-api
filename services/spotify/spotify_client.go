package spotify

import (
	"encoding/json"
	"io"
	"main/services/common"
	"math/rand"
	"net/http"
	"time"
)

func makeRequest[T any](logger *common.Logger, method string, url string, result *T) error {
	request, err := getRequest(logger, method, url)
	if err != nil {
		return err
	}

	client := &http.Client{}
	var response *http.Response

	attemptsLeft := 3
	for {
		if attemptsLeft == 0 {
			break
		}

		response, err = client.Do(request)
		if err != nil {
			return err
		}

		if response.StatusCode == http.StatusOK {
			break
		}

		if response.StatusCode == http.StatusTooManyRequests {
			timeout := getTimeout(attemptsLeft)
			time.Sleep(timeout)
		}

		attemptsLeft--
	}

	return getResponseContent(logger, response, &result)
}

func getRequest(logger *common.Logger, method string, url string) (*http.Request, error) {
	request, err := http.NewRequest(method, "https://api.spotify.com/v1/"+url, nil)
	if err != nil {
		return nil, err
	}

	accessToken, _ := getAccessToken(logger)
	if err != nil {
		logger.LogWarn(err.Error())
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

func getResponseContent[T any](logger *common.Logger, response *http.Response, result *T) error {
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logger.LogError(err, err.Error())
		return err
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logger.LogError(err, err.Error())
		return err
	}

	return nil
}

func getTimeout(attemptNumber int) time.Duration {
	base := timeouts[attemptNumber]
	jit := rand.Intn(jitInterval)

	return time.Duration(time.Millisecond * time.Duration(base+jit))
}

const jitInterval = 500

var timeouts = map[int]int{
	3: 500,
	2: 1000,
	1: 5000,
}
