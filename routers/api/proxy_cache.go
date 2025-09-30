package api

import (
    "sync"
    "time"
)

type cacheEntry struct {
    value      interface{}
    expiresAt  time.Time
}

type memoryCache struct {
    mu    sync.RWMutex
    items map[string]cacheEntry
}

var proxyCache = &memoryCache{items: make(map[string]cacheEntry)}

func (c *memoryCache) get(key string) (interface{}, bool) {
    c.mu.RLock()
    entry, ok := c.items[key]
    c.mu.RUnlock()
    if !ok {
        return nil, false
    }
    if time.Now().After(entry.expiresAt) {
        c.mu.Lock()
        delete(c.items, key)
        c.mu.Unlock()
        return nil, false
    }
    return entry.value, true
}

func (c *memoryCache) set(key string, val interface{}, ttl time.Duration) {
    c.mu.Lock()
    c.items[key] = cacheEntry{value: val, expiresAt: time.Now().Add(ttl)}
    c.mu.Unlock()
}

