package userdb

import (
	"testing"
	"time"
)

func TestCacheStringA(t *testing.T) {
	c := NewCache[string](time.Second, true)
	c.Set("foo", "bar", time.Second*2)
	time.Sleep(time.Second * 3)
	v, ok := c.Get("foo")
	if ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, nil, false)
	}
}

func TestCacheStringB(t *testing.T) {
	c := NewCache[string](time.Second, true)
	c.Set("foo", "bar", time.Second*2)
	v, ok := c.Get("foo")
	if !ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, "bar", true)
	}
	time.Sleep(time.Second * 3)
	v, ok = c.Get("foo")
	if ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, nil, false)
	}
}

func TestCacheStringC(t *testing.T) {
	c := NewCache[string](time.Second, true)
	c.Set("foo", "bar", time.Second*3)
	time.Sleep(time.Second * 2)
	v, ok := c.Get("foo")
	if !ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, "bar", true)
	}
	time.Sleep(time.Second * 2)
	v, ok = c.Get("foo")
	if !ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, "bar", true)
	}
	time.Sleep(time.Second * 4)
	v, ok = c.Get("foo")
	if ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, nil, false)
	}
}

func TestCacheStringD(t *testing.T) {
	c := NewCache[string](time.Second, false)
	c.Set("foo", "bar", time.Second*3)
	time.Sleep(time.Second * 2)
	v, ok := c.Get("foo")
	if !ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, "bar", true)
	}
	time.Sleep(time.Second * 2)
	v, ok = c.Get("foo")
	if ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, nil, false)
	}
	time.Sleep(time.Second * 4)
	v, ok = c.Get("foo")
	if ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, nil, false)
	}
}

func TestCacheStringE(t *testing.T) {
	c := NewCache[string](time.Second*10, false)
	c.Set("foo", "bar", time.Second*3)
	time.Sleep(time.Second * 2)
	v, ok := c.Get("foo")
	if !ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, "bar", true)
	}
	time.Sleep(time.Second * 2)
	v, ok = c.Get("foo")
	if ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, nil, false)
	}
	time.Sleep(time.Second * 4)
	v, ok = c.Get("foo")
	if ok {
		t.Errorf("cache.Get(foo) = %v, %v; want %v, %v", v, ok, nil, false)
	}
}

func TestCacheInt64(t *testing.T) {
	c := NewCache[int64](time.Second, true)
	c.Set(123, "bar", time.Second)
	v, ok := c.Get(123)
	if !ok {
		t.Errorf("cache.Get(123) = %v, %v; want %v, %v", v, ok, "bar", true)
	}
	time.Sleep(time.Second * 2)
	v, ok = c.Get(123)
	if ok {
		t.Errorf("cache.Get(123) = %v, %v; want %v, %v", v, ok, nil, false)
	}
}
