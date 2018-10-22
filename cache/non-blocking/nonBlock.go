package nonblockcache

// Non blocking Cache
// A non blocking cache has the ability to work on requests
// while waiting for the result of the function.
import (
	"sync"
)

// Cache is custom cache that uses map
// as the underlying type
type Cache struct {
	cache map[string]*entry
	sync.Mutex
}

type entry struct {
	res   result
	ready chan struct{}
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
	return &Cache{cache: make(map[string]*entry)}
}

// Get is used to retrieve the value of the cache from the Redis
func (c *Cache) Get(key string, f Func) ([]byte, error) {
	c.Lock()
	e := c.cache[key]
	if e == nil {
		// This is the first request for this key.
		// This goroutine becomes responsible for computing
		// the value and broadcasting the ready condition.
		e = &entry{ready: make(chan struct{})}
		c.cache[key] = e
		c.Unlock()
		e.res.value, e.res.err = f()
		// broadcast ready condition
		close(e.ready)
	} else {
		// This is a repeat request for this key
		c.Unlock()
		// If there was an existing entry (in the else block), its value
		// is not necessarily ready yet.
		// another goroutine could still be calling the slow function f.
		// This the calling goroutine must wait for the entry's "ready"
		// condition before it reads the entry's result.
		// It does this by reading a value from the ready channel, since
		// this operations blocks until the channel is closed.
		<-e.ready
	}
	// the variables e.res.value and e.res.err in the entry are shared among multiple goroutines
	// The goroutine that creates the entry sets their values.
	// and other goroutines read their values once the "ready" condition has been broadcast.
	// Despite being accessed by multiple goroutines, no mutex lock is necessary.
	// The closing of the ready channel happens before any other goroutine receives the broadcast event,
	// so the write to those variables in the first goroutine happens before they are read by subsequent goroutines.
	// Our concurrent, duplicate-suppressing, non-blocking cache is complete.
	return e.res.value, e.res.err
}
