package jyotish

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	data   JyotishData
	expiry time.Time
}

var responseCache sync.Map

const cacheTTL = 10 * time.Minute

func cacheKey(lat, lon float64, t time.Time) string {
	return fmt.Sprintf("%s/%.4f/%.4f", t.Format("2006-01-02T15:04"), lat, lon)
}

// CacheResponse stores a JyotishData response with a default TTL.
func CacheResponse(key string, data JyotishData) {
	responseCache.Store(key, cacheEntry{data: data, expiry: time.Now().Add(cacheTTL)})
}

// CachedResponse retrieves a cached JyotishData response.
func CachedResponse(key string) (JyotishData, bool) {
	val, ok := responseCache.Load(key)
	if !ok {
		return JyotishData{}, false
	}
	entry := val.(cacheEntry)
	if time.Now().After(entry.expiry) {
		responseCache.Delete(key)
		return JyotishData{}, false
	}
	return entry.data, true
}
