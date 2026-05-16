package jyotish

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	data   []byte
	expiry time.Time
}

var responseCache sync.Map

const cacheTTL = 10 * time.Minute

// CacheKey constructs a cache key from location and datetime strings.
func CacheKey(lat, lon float64, date, timeStr string) string {
	return fmt.Sprintf("%s/%s/%s/%s", date, timeStr, fmt.Sprintf("%.4f", lat), fmt.Sprintf("%.4f", lon))
}

// CacheResponse stores a raw response with a default TTL.
func CacheResponse(key string, data []byte) {
	responseCache.Store(key, cacheEntry{data: data, expiry: time.Now().Add(cacheTTL)})
}

// CachedResponse retrieves a cached raw response.
func CachedResponse(key string) ([]byte, bool) {
	val, ok := responseCache.Load(key)
	if !ok {
		return nil, false
	}
	entry := val.(cacheEntry)
	if time.Now().After(entry.expiry) {
		responseCache.Delete(key)
		return nil, false
	}
	return entry.data, true
}
