package group

import (
	"fmt"
	lru "gcache/Lru"
	"sync"
)

type cache struct {
	lru *lru.LruCache
	lk  sync.RWMutex

	maxbytes int
}

func newCache(maxbytes int) (c *cache) {
	return &cache{
		maxbytes: maxbytes,
	}
}

func (c *cache) add(key string, val lru.Value) {
	c.lk.Lock()
	defer c.lk.Unlock()

	if c.lru == nil {
		c.lru = lru.NewLruCache(c.maxbytes)
	}

	c.lru.Add(key, val)
}

func (c *cache) get(key string) (v lru.Value, err error) {
	c.lk.RLock()
	defer c.lk.RUnlock()

	if c.lru == nil {
		err = fmt.Errorf("[Error]: delay init...get before add")
		return
	}

	return c.lru.Get(key)
}
