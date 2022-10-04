package common

type SyncedResult[T any] struct {
	Result T
	Sync   string
}
