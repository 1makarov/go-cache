package cache

import (
	"fmt"
	"sync"
	"time"
)

const (
	errKeyBusy    = "key is busy"
	errValueEmpty = "value empty"
)

type Cache struct {
	s map[string]interface{}

	h *handler

	ch chan string

	m sync.Mutex
}

func New() *Cache {
	s := make(map[string]interface{})
	ch := make(chan string)
	h := initHandler(ch)

	c := &Cache{s: s, h: h, ch: ch}
	go c.run()

	return c
}

func (c *Cache) run() {
	for k := range c.ch {
		c.Delete(k)
	}
}

func (c *Cache) Set(k string, v interface{}) error {
	c.m.Lock()
	defer c.m.Unlock()
	_, ok := c.s[k]
	if ok {
		return fmt.Errorf(errKeyBusy)
	}
	c.s[k] = v

	return nil
}

func (c *Cache) SetWithExpire(k string, v interface{}, ttl time.Duration) error {
	t := time.Now().Add(ttl)

	if err := c.Set(k, v); err != nil {
		return err
	}
	c.h.add(k, t)

	return nil
}

func (c *Cache) Delete(k string) {
	c.m.Lock()
	delete(c.s, k)
	c.m.Unlock()
}

func (c *Cache) Get(k string) (interface{}, error) {
	v, ok := c.s[k]
	if !ok {
		return nil, fmt.Errorf(errValueEmpty)
	}
	return v, nil
}

func (c *Cache) Close() {
	close(c.ch)

	c.h.close()

	c.s = make(map[string]interface{})
}
