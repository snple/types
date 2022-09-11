package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/snple/types"
)

var (
	ErrNotFound = errors.New("cache: key not found")
)

type Value[T any] struct {
	Data    T
	TTL     time.Duration
	Updated time.Time
}

func newValue[T any](data T, ttl time.Duration) Value[T] {
	return Value[T]{
		Data:    data,
		TTL:     ttl,
		Updated: time.Now(),
	}
}

func (v *Value[T]) Alive() bool {
	return v != nil && (v.TTL == 0 || time.Since(v.Updated) <= v.TTL)
}

type Cache[T any] struct {
	data map[string]Value[T]
	lock sync.RWMutex
	miss func(key string) (T, time.Duration, error)
}

func NewCache[T any](miss func(key string) (T, time.Duration, error)) *Cache[T] {
	return &Cache[T]{
		data: make(map[string]Value[T]),
		lock: sync.RWMutex{},
		miss: miss,
	}
}

func (c *Cache[T]) GetValue(key string) types.Option[Value[T]] {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if value, ok := c.data[key]; ok {
		return types.Some(value)
	}

	return types.None[Value[T]]()
}

func (c *Cache[T]) Get(key string) types.Option[T] {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if value, ok := c.data[key]; ok && value.Alive() {
		return types.Some(value.Data)
	}

	return types.None[T]()
}

func (c *Cache[T]) GetWithMiss(key string) (types.Option[T], error) {
	if v := c.Get(key); v.IsSome() {
		return v, nil
	}

	if c.miss != nil {
		value, ttl, err := c.miss(key)
		if err != nil {
			return types.None[T](), err
		}

		c.Set(key, value, ttl)
		return types.Some(value), nil
	}

	return types.None[T](), ErrNotFound
}

func (c *Cache[T]) Set(key string, value T, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = newValue(value, ttl)
}

func (c *Cache[T]) AutoSet(key string, fn func(key string) (T, time.Duration, error), duration time.Duration) chan<- struct{} {
	quit := make(chan struct{})

	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if value, ttl, err := fn(key); err == nil {
					c.Set(key, value, ttl)
				}
			case <-quit:
				return
			}
		}
	}()

	return quit
}

func (c *Cache[T]) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, key)
}

func (c *Cache[T]) DeleteAll() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data = make(map[string]Value[T])
}

func (c *Cache[T]) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}

func (c *Cache[T]) GC() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key, value := range c.data {
		if !value.Alive() {
			delete(c.data, key)
		}
	}
}

func (c *Cache[T]) AutoGC(duration time.Duration) chan<- struct{} {
	quit := make(chan struct{})

	go func() {
		ticker := time.NewTicker(duration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.GC()
			case <-quit:
				return
			}
		}
	}()

	return quit
}
