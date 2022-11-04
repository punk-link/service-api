package helpers

func AlignSlice[T any](slice []T, divisor int) ([]T, []T) {
	var unalignedItems []T
	alignedItems := make([]T, 0)

	if len(slice) < divisor {
		unalignedItems = slice
		return alignedItems, unalignedItems
	}

	extraElements := len(slice) % divisor

	alignedItems = slice[0 : len(slice)-extraElements]
	unalignedItems = slice[len(slice)-extraElements:]

	return alignedItems, unalignedItems
}

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

func Distinct[T comparable](source []T) []T {
	distincted := make([]T, 0)
	distinctionMap := make(map[T]int, 0)
	for _, item := range source {
		if _, isDuplicated := distinctionMap[item]; !isDuplicated {
			distinctionMap[item] = 0
			distincted = append(distincted, item)
		}
	}

	return distincted
}

func getChunkEndPosition[T any](chunk []T, iterationStepEnd int) int {
	chunkEndPosition := len(chunk)

	if chunkEndPosition < iterationStepEnd {
		return chunkEndPosition
	}

	return iterationStepEnd
}
