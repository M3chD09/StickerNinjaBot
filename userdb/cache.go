package userdb

import (
	"sync"
	"time"
)

type key interface {
	int64 | string
}

type item struct {
	value  interface{}
	start  int64
	expire time.Duration
}

type cache[T key] struct {
	items            map[T]*item
	mux              *sync.RWMutex
	autoUpdateExpire bool
}

func NewCache[T key](tick time.Duration, autoUpdateExpire bool) *cache[T] {
	c := &cache[T]{
		items:            make(map[T]*item),
		mux:              &sync.RWMutex{},
		autoUpdateExpire: autoUpdateExpire,
	}
	go c.clean(tick)
	return c
}

func (c *cache[T]) Get(key T) (interface{}, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if v, ok := c.items[key]; !ok || time.Unix(v.start, 0).Add(v.expire).Unix() < time.Now().Unix() {
		delete(c.items, key)
		return nil, false
	} else {
		if c.autoUpdateExpire {
			v.start = time.Now().Unix()
		}
		return v.value, true
	}
}

func (c *cache[T]) Set(key T, value interface{}, expire time.Duration) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.items[key] = &item{
		value:  value,
		start:  time.Now().Unix(),
		expire: expire,
	}
}

func (c *cache[T]) Delete(key T) {
	c.mux.Lock()
	defer c.mux.Unlock()
	delete(c.items, key)
}

func (c *cache[T]) clean(tick time.Duration) {
	for range time.Tick(tick) {
		c.mux.Lock()
		for k, v := range c.items {
			if time.Unix(v.start, 0).Add(v.expire).Unix() < time.Now().Unix() {
				delete(c.items, k)
			}
		}
		c.mux.Unlock()
	}
}
