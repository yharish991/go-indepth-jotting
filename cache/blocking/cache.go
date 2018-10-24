package cache

import (
	"sync"
)

// Cache is custom cache that uses map
// as the underlying type
type Cache struct {
	cache map[string]*result
	sync.Mutex
}

type result struct {
	value []byte
	err   error
}

// Blocking Cache
// If the Cache request result in a miss, the cache must wait for the result
// of the slow function, until then it's is blocked.

// Func is used to compute the value
type Func func() ([]byte, error)

// NewCache is an Construtor
func NewCache() *Cache {
	return &Cache{cache: make(map[string]*result)}
}

// Get is used to retrieve the value of the cache from the Redis
func (c *Cache) Get(key string, f Func) ([]byte, error) {
	// since cache is the critical section of the data,
	// and with the whole function being the monitor based
	// synchronization.
	// but this blocking will result in the duplicate value
	// Imaging two goroutine accessing the cache GET at the same time.
	c.Lock()
	res, ok := c.cache[key]
	c.Unlock()
	if !ok {
		res = &result{}
		res.value, res.err = f()
		c.Lock()
		c.cache[key] = res
		c.Unlock()
	}
	return res.value, res.err
}
