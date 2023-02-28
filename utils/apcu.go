package utils

import (
	"fmt"
	"sync"
	"time"
)

// Create a global sync.Map instance for storing cache entries
var cache = &sync.Map{}

type CacheEntry struct {
	available      bool
	expirationTime time.Time
}

func PrintCacheValues() {
	cache.Range(func(key, value interface{}) bool {
		// Convert the value to cacheEntry type
		entry, ok := value.(*CacheEntry)
		if !ok {
			return true // skip this entry
		}

		// Print the key and value
		fmt.Printf("Key: %v, Value: %v %v\n", key, entry.available, entry.expirationTime)

		return true // continue iterating
	})
}

func ApcuExists(key string) bool {
	// Check if the given key exists in the sync.Map
	PrintCacheValues()
	_, ok := cache.Load(key)
	return ok
}

func ApcuFetch(key string) *CacheEntry {
	// Load the cache entry from the sync.Map
	entry, _ := cache.Load(key)

	// Check if the cache entry has expired
	if time.Now().After(entry.(*CacheEntry).expirationTime) {
		cache.Delete(key)
	}

	return entry.(*CacheEntry)
}
