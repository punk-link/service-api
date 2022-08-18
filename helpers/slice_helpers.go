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

func getChunkEndPosition[T any](chunk []T, iterationStepEnd int) int {
	chunkEndPosition := len(chunk)

	if chunkEndPosition < iterationStepEnd {
		return chunkEndPosition
	}

	return iterationStepEnd
}
