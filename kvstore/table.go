// Package KVstore ;
//
// Inspired from: https://github.com/cheshir/ttlcache
// And from: https://github.com/zekroTJA/timedmap
//
package kvstore

import (
	"encoding/json"
	"sync"
	"time"
)

type CacheTable struct {
	sync.RWMutex
	items        map[string]item
	done         chan bool
	cleanRunning bool
}

type item struct {
	expire int64 // Unix micro
	value  interface{}
}

func NewCache() *CacheTable {
	return &CacheTable{
		items: make(map[string]item),
		done:  make(chan bool),
	}
}

func (c *CacheTable) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

// internal get + delete expired items (probably a bad idea)
func (c *CacheTable) get(key string) (item, bool) {
	c.RLock()
	cacheItem, ok := c.items[key]
	c.RUnlock()
	if !ok {
		return item{}, false
	}
	if cacheItem.expire > 0 && cacheItem.expire < time.Now().UnixMicro() {
		c.Lock()
		delete(c.items, key)
		c.Unlock()
		return item{}, false
	}
	return cacheItem, true
}

func (c *CacheTable) Exists(key string) bool {
	_, ok := c.get(key)
	return ok
}

// Get returns stored record.
// First returned: the stored value.
// Second returned: existence flag like in the map.
func (c *CacheTable) Get(key string) (interface{}, bool) {
	cacheItem, ok := c.get(key)
	return cacheItem.value, ok
}

// Set adds record in the cache with given ttl.
// If TTL is less than zero, it will be stored forever.
func (c *CacheTable) Set(key string, value interface{}, ttl time.Duration) {
	cacheItem := item{value: value}
	if ttl == 0 {
		cacheItem.expire = 0
	} else if ttl < 0 {
		cacheItem.expire = -1
	} else {
		cacheItem.expire = time.Now().Add(ttl).UnixMicro()
	}
	c.Lock()
	c.items[key] = cacheItem
	c.Unlock()
}

// Delete deletes the key and its value from the cache.
func (c *CacheTable) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

// Count returns how many items are currently stored in the cache.
func (c *CacheTable) Count() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.items)
}

// cleanUp removes outdated items from the cache.
func (c *CacheTable) cleanUp() {
	now := time.Now().UnixMicro()
	c.Lock()
	defer c.Unlock()

	for key, item := range c.items {
		if item.expire > 0 && item.expire < now {
			delete(c.items, key)
		}
	}
}

// cleanupLoop blocks the loop executing the cleanup.
func (c *CacheTable) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			c.cleanUp()
		case <-c.done:
			ticker.Stop()
			return
		}
	}
}

// StartCleaner starts the cleanup loop controlled
// by an internal ticker with the given interval.
// If the cleanup loop is already running, it will be
// stopped and restarted using the new specification.
func (c *CacheTable) StartCleaner(interval time.Duration) {
	if c.cleanRunning {
		c.StopCleaner()
	}
	go c.cleanupLoop(interval)
}

// StopCleaner stops the cleaner go routine and timer.
// This should always be called after exiting a scope
// where CacheTable is used that the data can be cleaned
// up correctly.
func (c *CacheTable) StopCleaner() {
	if !c.cleanRunning {
		return
	}
	c.done <- true
}
