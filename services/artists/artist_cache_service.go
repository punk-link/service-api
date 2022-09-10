package artists

import (
	"main/models/artists"
	"sync"
	"time"
)

type cacheEntry struct {
	Treshold time.Time
	Value    artists.Artist
}

type ArtistCacheService struct {
	cache map[string]cacheEntry
	mutex sync.Mutex
}

func ConstructArtistCacheService() *ArtistCacheService {
	cache := make(map[string]cacheEntry)

	service := ArtistCacheService{
		cache: cache,
	}

	go service.watch()

	return &service
}

func (t *ArtistCacheService) TryGet(key string) (artists.Artist, bool) {
	if key == "" {
		return artists.Artist{}, false
	}

	if entry, isCached := t.cache[key]; isCached {
		return entry.Value, true
	}

	return artists.Artist{}, false
}

func (t *ArtistCacheService) Remove(key string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	delete(t.cache, key)
}

func (t *ArtistCacheService) Set(key string, value artists.Artist, interval time.Duration) {
	if key == "" {
		return
	}

	treshold := time.Now().UTC().Add(interval)
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.cache[key] = cacheEntry{
		Treshold: treshold,
		Value:    value,
	}
}

func (t *ArtistCacheService) watch() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		for key, value := range t.cache {
			if time.Now().UTC().After(value.Treshold) {
				t.Remove(key)
			}
		}
	}
}
