package go_cache

import "sync"

type Cache struct {
	storage map[string]interface{}

	s sync.Mutex
}

func New() *Cache {
	s := make(map[string]interface{})
	return &Cache{storage: s}
}

func (c *Cache) Set(k string, v interface{}) {
	c.s.Lock()
	c.storage[k] = v
	c.s.Unlock()
	return
}

func (c *Cache) Delete(k string) {
	c.s.Lock()
	delete(c.storage, k)
	c.s.Unlock()
	return
}

func (c *Cache) Get(k string) interface{} {
	return c.storage[k]
}
