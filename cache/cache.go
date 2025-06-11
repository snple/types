package cache

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/snple/types"
)

var (
	ErrNotFound = errors.New("cache: key not found")
)

type Value[V any] struct {
	Data    V
	TTL     time.Duration
	Updated time.Time
}

func NewValue[V any](data V, ttl time.Duration) Value[V] {
	return Value[V]{
		Data:    data,
		TTL:     ttl,
		Updated: time.Now(),
	}
}

func (v *Value[V]) Alive() bool {
	return v != nil && (v.TTL == 0 || time.Since(v.Updated) <= v.TTL)
}

type Cache[K comparable, V any] struct {
	data map[K]Value[V]
	lock sync.RWMutex
	miss func(ctx context.Context, key K) (V, time.Duration, error)
}

func NewCache[K comparable, V any](miss func(ctx context.Context, key K) (V, time.Duration, error)) *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]Value[V]),
		lock: sync.RWMutex{},
		miss: miss,
	}
}

func (c *Cache[K, V]) GetValue(key K) types.Option[Value[V]] {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if value, ok := c.data[key]; ok {
		return types.Some(value)
	}

	return types.None[Value[V]]()
}

func (c *Cache[K, V]) Get(key K) types.Option[V] {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if value, ok := c.data[key]; ok && value.Alive() {
		return types.Some(value.Data)
	}
	return types.None[V]()
}

func (c *Cache[K, V]) GetWithMiss(ctx context.Context, key K) (types.Option[V], error) {
	if v := c.Get(key); v.IsSome() {
		return v, nil
	}

	if c.miss != nil {
		value, ttl, err := c.miss(ctx, key)
		if err != nil {
			return types.None[V](), err
		}

		c.Set(key, value, ttl)
		return types.Some(value), nil
	}

	return types.None[V](), ErrNotFound
}

func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[key] = NewValue(value, ttl)
}

func (c *Cache[K, V]) AutoSet(key K, fn func(key K) (V, time.Duration, error), duration time.Duration) chan<- struct{} {
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

func (c *Cache[K, V]) Delete(key K) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.data, key)
}

func (c *Cache[K, V]) DeleteAll() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data = make(map[K]Value[V])
}

func (c *Cache[K, V]) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}

func (c *Cache[K, V]) GC() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key, value := range c.data {
		if !value.Alive() {
			delete(c.data, key)
		}
	}
}

func (c *Cache[K, V]) AutoGC(duration time.Duration) chan<- struct{} {
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
