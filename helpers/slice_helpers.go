package helpers

import dataStructures "main/data-structures"

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
	distinctionSet := dataStructures.MakeHashSet([]T{})
	for _, item := range source {
		if !distinctionSet.Contains(item) {
			distinctionSet.Add(item)
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
