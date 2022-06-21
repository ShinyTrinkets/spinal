package kvstore

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

// Store returns the existing cache table with given name or
// creates a new one if the table doesn't exist yet.
func Store(table string) *CacheTable {
	mutex.RLock()
	t, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		t, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			cache[table] = &CacheTable{name: table, done: make(chan bool)}
		}
		mutex.Unlock()
	}

	return t
}
