package dataStructures

type HashSet[T any] struct {
	innerMap map[any]bool
}

func MakeEmptyHashSet[T any]() HashSet[T] {
	return HashSet[T](MakeHashSet([]T{}))
}

func MakeHashSet[T any](values []T) HashSet[T] {
	innerMap := make(map[any]bool, len(values))
	for _, value := range values {
		innerMap[value] = true
	}

	return HashSet[T]{
		innerMap: innerMap,
	}
}

func (t *HashSet[T]) Add(value T) {
	t.innerMap[value] = true
}

func (t *HashSet[T]) AsMap() map[any]bool {
	return t.innerMap
}

func (t *HashSet[T]) Contains(value T) bool {
	return t.innerMap[value]
}
