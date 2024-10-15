package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
type Cache struct {
	cachedMap map[string]cacheEntry
	mut       sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{cachedMap: map[string]cacheEntry{}, mut: sync.Mutex{}}
	c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) error {
	c.mut.Lock()
	defer c.mut.Unlock()
	if len(val) == 0 {
		return fmt.Errorf("error: array's length is 0")
	}
	c.cachedMap[key] = cacheEntry{createdAt: time.Now(), val: val}
	return nil
}
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mut.Lock()
	defer c.mut.Unlock()
	entry, ok := c.cachedMap[key]
	if !ok {
		return []byte{}, ok
	}
	return entry.val, ok
}
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.mut.Lock()
			for key, val := range c.cachedMap {
				if time.Since(val.createdAt) > interval {
					delete(c.cachedMap, key)
				}
			}
			c.mut.Unlock()
		}
	}()
}
