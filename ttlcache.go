// ttlcache.go
package gomemorycache

import (
	"sync"
	"time"
)

type item[V any] struct {
	value  V
	expiry time.Time
}

func (i item[V]) isExpired() bool {
	return time.Now().After(i.expiry)
}

type TTLCache[K comparable, V any] struct {
	items map[K]item[V]
	mutex sync.Mutex
}

func NewTTLCache[K comparable, V any]() *TTLCache[K, V] {
	c := &TTLCache[K, V]{
		items: make(map[K]item[V]),
	}
	go func() {
		for range time.Tick(1 * time.Minute) {
			c.mutex.Lock()
			for key, item := range c.items {
				if item.isExpired() {
					delete(c.items, key)
				}
			}
			c.mutex.Unlock()
		}
	}()
	return c
}

func (c *TTLCache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = item[V]{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

func (c *TTLCache[K, V]) Get(key K) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.items[key]
	if !found {
		return item.value, false
	}
	if item.isExpired() {
		delete(c.items, key)
		return item.value, false
	}
	return item.value, true
}

func (c *TTLCache[K, V]) Remove(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

func (c *TTLCache[K, V]) Pop(key K) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.items[key]
	if !found {
		return item.value, false
	}
	delete(c.items, key)
	if item.isExpired() {
		return item.value, false
	}
	return item.value, true
}
