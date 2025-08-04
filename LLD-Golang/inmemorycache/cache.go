package main

import (
	"sync"
	"time"
)

// value and expiration time
type Entry struct {
	Value    interface{}
	ExpireAt time.Time
}

// interface for eviction strategies
type EvictionPolicy interface {
	Access(key string)
	Add(key string)
	Evict() string
	Remove(key string)
}

// cache storage with eviction policy and TTL
type Cache struct {
	capacity int
	storage  map[string]*Entry
	policy   EvictionPolicy
	mu       sync.Mutex
}

func NewCache(capacity int, policy EvictionPolicy) *Cache {
	return &Cache{
		capacity: capacity,
		storage:  make(map[string]*Entry),
		policy:   policy,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.storage[key]
	if !ok {
		return nil, false
	}

	if !entry.ExpireAt.IsZero() && time.Now().After(entry.ExpireAt) {
		delete(c.storage, key)
		c.policy.Remove(key)
		return nil, false
	}

	c.policy.Access(key)

	return entry.Value, true
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}

	if entry, ok := c.storage[key]; ok {
		entry.Value = value
		entry.ExpireAt = expireAt
		c.policy.Access(key)
		return
	}

	// New key: evict if full
	if len(c.storage) >= c.capacity {
		evicted := c.policy.Evict()
		if evicted != "" {
			delete(c.storage, evicted)
		}
	}

	c.storage[key] = &Entry{Value: value, ExpireAt: expireAt}
	c.policy.Add(key)
}

// Delete removes the given key from the cache, and from the underlying
// eviction policy.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.storage[key]; ok {
		delete(c.storage, key)
		c.policy.Remove(key)
	}
}
