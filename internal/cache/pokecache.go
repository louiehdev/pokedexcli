package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu        sync.Mutex
	PokeCache map[string]CacheEntry
}

type CacheEntry struct {
	CreatedAt time.Time
	Data      []byte
}

func (c *Cache) Add(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.PokeCache[key] = CacheEntry{CreatedAt: time.Now(), Data: data}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, exists := c.PokeCache[key]
	if exists {
		return entry.Data, true
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.PokeCache {
			if time.Since(entry.CreatedAt) > interval {
				delete(c.PokeCache, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{PokeCache: make(map[string]CacheEntry)}
	go newCache.reapLoop(interval)
	return newCache
}
