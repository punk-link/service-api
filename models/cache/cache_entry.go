package cache

import "time"

type CacheEntry struct {
	Treshold time.Time
	Value    interface{}
}
