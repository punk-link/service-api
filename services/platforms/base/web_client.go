package base

import (
	"encoding/json"
	"io"
	"main/helpers"
	commonModels "main/models/common"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/punk-link/logger"
)

func MakeBatchRequestWithSync[T any](logger *logger.Logger, syncedHttpRequests []commonModels.SyncedHttpRequest) []commonModels.SyncedResult[T] {
	iterationStep := runtime.NumCPU()
	alignedRequests, unalignedRequests := helpers.AlignSlice(syncedHttpRequests, iterationStep)

	var wg sync.WaitGroup
	chanResults := make(chan commonModels.SyncedResult[T])

	for i := 0; i < len(alignedRequests); i = i + iterationStep {
		wg.Add(iterationStep)

		for j := 0; j < iterationStep; j++ {
			go processRequestAsync(&wg, chanResults, logger, alignedRequests[i+j])
		}

		time.Sleep(REQUEST_BATCH_TIMEOUT_DURATION_MILLISECONDS)
	}

	for _, request := range unalignedRequests {
		wg.Add(1)
		go processRequestAsync(&wg, chanResults, logger, request)
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

func MakeRequest[T any](logger *logger.Logger, request *http.Request, result *T) error {
	client := &http.Client{}
	var response *http.Response

	var err error
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
			logger.LogWarn("Web request ends with a status code %v", response.StatusCode)
			time.Sleep(getTimeoutDuration(attemptsLeft))
		}

		attemptsLeft--
	}

	return unmarshalResponseContent(logger, response, &result)
}

func processRequestAsync[T any](wg *sync.WaitGroup, results chan<- commonModels.SyncedResult[T], logger *logger.Logger, syncedHttpRequest commonModels.SyncedHttpRequest) {
	defer wg.Done()

	var responseContent T
	err := MakeRequest(logger, syncedHttpRequest.HttpRequest, &responseContent)
	if err != nil {
		return
	}

	results <- commonModels.SyncedResult[T]{
		Result: responseContent,
		Sync:   syncedHttpRequest.Sync,
	}
}

func getTimeoutDuration(attemptNumber int) time.Duration {
	base := _timeoutValues[attemptNumber]
	jit := rand.Intn(JITTER_INTERVAL_MILLISECONDS)

	return time.Millisecond * time.Duration(base+jit)
}

func unmarshalResponseContent[T any](logger *logger.Logger, response *http.Response, result *T) error {
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

var _timeoutValues = map[int]int{
	3: 500,
	2: 1000,
	1: 5000,
}

const JITTER_INTERVAL_MILLISECONDS = 500
const REQUEST_BATCH_TIMEOUT_DURATION_MILLISECONDS = time.Millisecond * 100
