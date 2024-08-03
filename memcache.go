// memcache.go
package gomemorycache

import "sync"

type Cache[K comparable, V any] struct {
	items map[K]V
	mutex sync.Mutex
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]V),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, found := c.items[key]
	return value, found
}

func (c *Cache[K, V]) Remove(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.items, key)
}

func (c *Cache[K, V]) Pop(key K) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, found := c.items[key]
	if found {
		delete(c.items, key)
	}
	return value, found
}
