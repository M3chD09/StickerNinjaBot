package userdb

import (
	"sync"
	"time"
)

type key interface {
	int64 | string
}

type item struct {
	value      interface{}
	start      int64
	expiration time.Duration
}

type Cache[T key] struct {
	items                map[T]*item
	mux                  *sync.RWMutex
	autoUpdateExpiration bool
}

func NewCache[T key](tick time.Duration, autoUpdateExpiration bool) *Cache[T] {
	c := &Cache[T]{
		items:                make(map[T]*item),
		mux:                  &sync.RWMutex{},
		autoUpdateExpiration: autoUpdateExpiration,
	}
	go c.clean(tick)
	return c
}

func (c *Cache[T]) Get(key T) (interface{}, bool) {
	if c.autoUpdateExpiration {
		// autoUpdateExpiration writes v.start, so it must hold the write
		// lock to avoid a data race with clean() and other Get calls.
		c.mux.Lock()
		defer c.mux.Unlock()
		if v, ok := c.items[key]; !ok || time.Unix(v.start, 0).Add(v.expiration).Unix() < time.Now().Unix() {
			delete(c.items, key)
			return nil, false
		} else {
			v.start = time.Now().Unix()
			return v.value, true
		}
	}

	c.mux.RLock()
	defer c.mux.RUnlock()
	if v, ok := c.items[key]; !ok || time.Unix(v.start, 0).Add(v.expiration).Unix() < time.Now().Unix() {
		return nil, false
	} else {
		return v.value, true
	}
}

func (c *Cache[T]) Set(key T, value interface{}, expiration time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.items[key] = &item{
		value:      value,
		start:      time.Now().Unix(),
		expiration: expiration,
	}
}

func (c *Cache[T]) Delete(key T) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.items, key)
}

func (c *Cache[T]) clean(tick time.Duration) {
	for range time.Tick(tick) {
		c.mux.Lock()
		for k, v := range c.items {
			if time.Unix(v.start, 0).Add(v.expiration).Unix() < time.Now().Unix() {
				delete(c.items, k)
			}
		}
		c.mux.Unlock()
	}
}
