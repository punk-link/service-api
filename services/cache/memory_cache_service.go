package cache

import (
	"sync"
	"time"

	"main/models/cache"
)

type MemoryCacheService struct {
	cache map[string]cache.CacheEntry
	mutex sync.Mutex
}

func ConstructMemoryCacheService() *MemoryCacheService {
	cache := make(map[string]cache.CacheEntry)

	service := MemoryCacheService{
		cache: cache,
	}

	go service.watch()

	return &service
}

func (t *MemoryCacheService) Remove(key string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	delete(t.cache, key)
}

func (t *MemoryCacheService) Set(key string, value interface{}, interval time.Duration) {
	if key == "" {
		return
	}

	treshold := time.Now().UTC().Add(interval)
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.cache[key] = cache.CacheEntry{
		Treshold: treshold,
		Value:    value,
	}
}

func (t *MemoryCacheService) TryGet(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	if entry, isCached := t.cache[key]; isCached {
		return entry.Value, true
	}

	return nil, false
}

func (t *MemoryCacheService) watch() {
	ticker := time.NewTicker(LIFETIME_VALIDATION_INTERVAL)
	for range ticker.C {
		for key, value := range t.cache {
			if time.Now().UTC().After(value.Treshold) {
				t.Remove(key)
			}
		}
	}
}

const LIFETIME_VALIDATION_INTERVAL time.Duration = 5 * time.Second
