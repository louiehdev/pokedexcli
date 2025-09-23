package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	PokeCache map[string]CacheEntry
}

type CacheEntry struct {
	CreatedAt time.Time
	Data []byte
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

	for t := range ticker.C {
		c.mu.Lock()
		for key, entry := range c.PokeCache {
			difference := t.Sub(entry.CreatedAt)
			if difference > interval {
				delete(c.PokeCache, key)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval int) *Cache {
	duration := time.Duration(interval) * time.Second
	newCache := &Cache{PokeCache: make(map[string]CacheEntry)}
	go newCache.reapLoop(duration)
	return newCache
}

