package shared

import (
	"testing"
	"time"
)

// clearCache resets the cache. This helps ensure that tests are isolated.
func clearCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cache = make(map[string]cacheItem)
}

func TestGetCacheMiss(t *testing.T) {
	clearCache()
	// For a non-existent key, we expect no data and false.
	value, ok := GetCache("nonexistent")
	if ok || value != "" {
		t.Errorf("expected miss for key 'nonexistent', got value=%q, ok=%v", value, ok)
	}
}

func TestGetCacheHit(t *testing.T) {
	clearCache()
	// Set a value with a TTL of 5 seconds.
	SetCache("testKey", "testValue", 5)
	value, ok := GetCache("testKey")
	if !ok || value != "testValue" {
		t.Errorf("expected hit for key 'testKey' with value 'testValue', got value=%q, ok=%v", value, ok)
	}
}

func TestGetCacheExpired(t *testing.T) {
	clearCache()
	// Set a value with a TTL of 1 second.
	SetCache("expireKey", "expireValue", 1)
	// Wait for the key to expire.
	time.Sleep(2 * time.Second)
	value, ok := GetCache("expireKey")
	if ok || value != "" {
		t.Errorf("expected expired key 'expireKey' to return empty value and false, got value=%q, ok=%v", value, ok)
	}
}
