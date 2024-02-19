package golrucache

import (
	"sync"

	"example.com/GoLRUCache/lru"
)

// encapsulation  of lru with mutex

type cache struct {
	mutex      sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// lazy initialization
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)

}

func (c *cache) get(key string) (value ByteView, exist bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		return
	}

	if value, exist := c.lru.Get(key); exist {
		return value.(ByteView), exist
	}
	return
}
