package kvstore

import (
	"sync"
)

var (
	cache = make(map[string]*CacheTable)
	mutex sync.RWMutex
)

func List() []string {
	list := []string{}
	mutex.RLock()
	for k := range cache {
		list = append(list, k)
	}
	mutex.RUnlock()
	return list
}

// Store returns the existing cache table with given name or
// creates a new one if the table doesn't exist yet.
func Store(table string) *CacheTable {
	mutex.RLock()
	_, ok := cache[table]
	mutex.RUnlock()

	if !ok {
		mutex.Lock()
		_, ok = cache[table]
		// Double check whether the table exists or not.
		if !ok {
			cache[table] = &CacheTable{name: table, done: make(chan bool)}
		}
		mutex.Unlock()
	}

	mutex.RLock()
	defer mutex.RUnlock()
	t := cache[table]
	return t
}
