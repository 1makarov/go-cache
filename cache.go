package cache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	tick *time.Ticker

	close bool

	values sync.Map
}

const (
	minInterval = time.Nanosecond

	errEmptyValue  = "empty value"
	errKeyIsBusy   = "key is busy"
	errKeyNotFound = "key not found"
)

func New() *Cache {
	return NewWithInterval(minInterval)
}

func NewWithInterval(interval time.Duration) *Cache {
	c := &Cache{
		tick: time.NewTicker(interval),
	}

	go c.waiter()

	return c
}

func (c *Cache) Close() {
	c.close = true
}

func (c *Cache) ChangeInterval(interval time.Duration) {
	c.tick.Reset(interval)
}

func (c *Cache) waiter() {
	defer c.tick.Stop()

	for range c.tick.C {
		if c.close {
			return
		}

		c.handler()
	}
}

type value struct {
	data    interface{}
	expires int64
}

func (c *Cache) handler() {
	now := time.Now().UnixNano()

	c.values.Range(func(k, va interface{}) bool {
		v := va.(value)

		if v.expires > 0 && now >= v.expires {
			c.values.Delete(k)
		}

		return true
	})
}

func (c *Cache) Get(k interface{}) (interface{}, error) {
	if v, ok := c.values.Load(k); ok {
		return v.(value).data, nil
	}
	return nil, fmt.Errorf(errEmptyValue)
}

func (c *Cache) GetAndDelete(k interface{}) (interface{}, error) {
	if v, ok := c.values.LoadAndDelete(k); ok {
		return v.(value).data, nil
	}
	return nil, fmt.Errorf(errEmptyValue)
}

func (c *Cache) Set(k, v interface{}) error {
	if _, ok := c.values.LoadOrStore(k, value{data: v}); !ok {
		return nil
	}
	return fmt.Errorf(errKeyIsBusy)
}

func (c *Cache) SetWithDuration(k, v interface{}, d time.Duration) error {
	if _, ok := c.values.LoadOrStore(k, value{data: v, expires: time.Now().Add(d).UnixNano()}); !ok {
		return nil
	}
	return fmt.Errorf(errKeyIsBusy)
}

func (c *Cache) Delete(k interface{}) error {
	if _, ok := c.values.LoadAndDelete(k); ok {
		return nil
	}
	return fmt.Errorf(errKeyNotFound)
}
