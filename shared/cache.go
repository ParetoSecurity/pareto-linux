package shared

import (
	"sync"

	"time"
)

type cacheItem struct {
	data string

	timestamp time.Time
}

var (
	cache = make(map[string]cacheItem)

	cacheMutex sync.RWMutex
)

func GetCache(key string) (string, bool) {

	cacheMutex.RLock()

	defer cacheMutex.RUnlock()

	if item, exists := cache[key]; exists {

		return item.data, true

	}

	return "", false

}

func SetCache(key string, value string, ttlSeconds int) {

	cacheMutex.Lock()

	defer cacheMutex.Unlock()

	cache[key] = cacheItem{

		data: value,

		timestamp: time.Now(),
	}

}
