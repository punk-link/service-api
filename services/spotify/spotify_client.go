package spotify

import (
	"encoding/json"
	"io"
	"main/helpers"
	commonModels "main/models/common"
	"main/services/common"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func makeBatchRequestWithSync[T any](logger *common.Logger, method string, syncedUrls []commonModels.SyncedUrl) []commonModels.SyncedResult[T] {
	iterationStep := runtime.NumCPU()

	mainLoop, reducedLoop := helpers.DivideChunkToLoops(syncedUrls, iterationStep)

	var wg sync.WaitGroup
	chanResults := make(chan commonModels.SyncedResult[T])

	for i := 0; i < len(mainLoop); i = i + iterationStep {
		wg.Add(iterationStep)

		for j := 0; j < iterationStep; j++ {
			go makeBatchRequestInternal(&wg, chanResults, logger, method, mainLoop[i+j])
		}

		time.Sleep(REQUEST_BATCH_TIMEOUT_DURATION_MILLISECONDS)
	}

	for _, chunk := range reducedLoop {
		wg.Add(1)
		go makeBatchRequestInternal(&wg, chanResults, logger, method, chunk)
	}

	go func() {
		wg.Wait()
		close(chanResults)
	}()

	syncedResults := make([]commonModels.SyncedResult[T], 0)
	for result := range chanResults {
		syncedResults = append(syncedResults, result)
	}

	return syncedResults
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
			logger.LogError(err, err.Error())
			return err
		}

		if response.StatusCode == http.StatusOK {
			break
		}

		if response.StatusCode == http.StatusTooManyRequests {
			logger.LogWarn("Spotify request ends with a status code %v", response.StatusCode)

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
	jit := rand.Intn(JITTER_INTERVAL_MILLISECONDS)

	return time.Duration(time.Millisecond * time.Duration(base+jit))
}

func makeBatchRequestInternal[T any](wg *sync.WaitGroup, results chan<- commonModels.SyncedResult[T], logger *common.Logger, method string, syncedUrls commonModels.SyncedUrl) {
	defer wg.Done()

	var responseContent T
	_ = makeRequest(logger, method, syncedUrls.Url, &responseContent)

	results <- commonModels.SyncedResult[T]{
		Result: responseContent,
		Sync:   syncedUrls.Sync,
	}
}

var timeouts = map[int]int{
	3: 500,
	2: 1000,
	1: 5000,
}

const JITTER_INTERVAL_MILLISECONDS = 500
const REQUEST_BATCH_TIMEOUT_DURATION_MILLISECONDS = time.Millisecond * 100
