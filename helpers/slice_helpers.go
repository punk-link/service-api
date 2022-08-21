package helpers

func Chunk[T any](source []T, chunkSize int) [][]T {
	chunkCount := (len(source) / chunkSize) + 1

	chunked := make([][]T, chunkCount)
	skip := 0
	for i := 0; i < chunkCount; i++ {
		values := source[skip:getChunkEndPosition(source, skip+chunkSize)]

		chunked[i] = values
		skip += chunkSize
	}

	return chunked
}

func DivideChunkToLoops[T any](chunkedUrls []T, iterationStep int) ([]T, []T) {
	var reducedLoop []T
	mainLoop := make([]T, 0)

	if len(chunkedUrls) < iterationStep {
		reducedLoop = chunkedUrls
	} else {
		extraElements := len(chunkedUrls) % iterationStep

		mainLoop = chunkedUrls[0 : len(chunkedUrls)-extraElements]
		reducedLoop = chunkedUrls[len(chunkedUrls)-extraElements:]
	}

	return mainLoop, reducedLoop
}

func getChunkEndPosition[T any](chunk []T, iterationStepEnd int) int {
	chunkEndPosition := len(chunk)

	if chunkEndPosition < iterationStepEnd {
		return chunkEndPosition
	}

	return iterationStepEnd
}
