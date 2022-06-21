// Package KVstore ;
//
// Inspired from: https://github.com/cheshir/ttlcache
// And from: https://github.com/zekroTJA/timedmap
//
package kvstore

import (
	"sync"
	"time"
)

type CacheTable struct {
	name         string
	items        sync.Map
	done         chan bool
	cleanRunning bool
}

type item struct {
	expire int64 // Unix micro
	value  interface{}
}

// internal get + delete expired items (probably a bad idea)
func (table *CacheTable) get(key string) (item, bool) {
	x, ok := table.items.Load(key)
	if !ok {
		return item{}, false
	}
	cacheItem := x.(item)
	if cacheItem.expire > 0 && cacheItem.expire < time.Now().UnixMicro() {
		table.items.Delete(key)
		return item{}, false
	}
	return cacheItem, true
}

func (table *CacheTable) Exists(key string) bool {
	_, ok := table.get(key)
	return ok
}

// Get returns stored record.
// First returned: the stored value.
// Second returned: existence flag like in the map.
func (table *CacheTable) Get(key string) (interface{}, bool) {
	cacheItem, ok := table.get(key)
	return cacheItem.value, ok
}

// Set adds record in the cache with given ttl.
// If TTL is less than zero, it will be stored forever.
func (table *CacheTable) Set(key string, value interface{}, ttl time.Duration) {
	cacheItem := item{
		expire: time.Now().Add(ttl).UnixMicro(),
		value:  value,
	}
	if cacheItem.expire < 0 {
		cacheItem.expire = -1
	}
	table.items.Store(key, cacheItem)
}

// Delete deletes the key and its value from the cache.
func (table *CacheTable) Delete(key string) {
	table.items.Delete(key)
}

// Count returns how many items are currently stored in the cache.
func (table *CacheTable) Count() int {
	var i int
	table.items.Range(func(k, x interface{}) bool {
		i++
		return true
	})
	return i
}

// cleanUp removes outdated items from the cache.
func (table *CacheTable) cleanUp() {
	now := time.Now().UnixMicro()
	table.items.Range(func(key, x interface{}) bool {
		cacheItem := x.(item)
		if cacheItem.expire > 0 && cacheItem.expire < now {
			table.items.Delete(key)
		}
		return true
	})
}

// cleanupLoop blocks the loop executing the cleanup.
func (table *CacheTable) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			table.cleanUp()
		case <-table.done:
			ticker.Stop()
			return
		}
	}
}

// StartCleaner starts the cleanup loop controlled
// by an internal ticker with the given interval.
// If the cleanup loop is already running, it will be
// stopped and restarted using the new specification.
func (table *CacheTable) StartCleaner(interval time.Duration) {
	if table.cleanRunning {
		table.StopCleaner()
	}
	go table.cleanupLoop(interval)
}

// StopCleaner stops the cleaner go routine and timer.
// This should always be called after exiting a scope
// where CacheTable is used that the data can be cleaned
// up correctly.
func (table *CacheTable) StopCleaner() {
	if !table.cleanRunning {
		return
	}
	table.done <- true
}
