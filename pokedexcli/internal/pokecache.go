package internal

import (
    "sync"
    "time"
)

type cacheEntry struct {
    createdAt time.Time
    val       []byte
}

type Cache struct {
    cache map[string]cacheEntry
    mu       *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
    c := Cache{
        cache: make(map[string]cacheEntry),
        mu:    &sync.Mutex{},
    }
    go c.reapLoop(interval)
    return c
}

func (c Cache) Add(key string, val []byte) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.cache[key] = cacheEntry{createdAt: time.Now() ,val: val}
}

func (c Cache) Get(key string) ([]byte, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if entry, found := c.cache[key]; found {
        return entry.val, found
    }
    return nil, false
}

func (c Cache) reapLoop(interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            c.mu.Lock()
            for key, entry := range c.cache {
                elapsed := time.Since(entry.createdAt)
                if elapsed > interval {
                    delete(c.cache, key)
                }
            }
            c.mu.Unlock()
        }
    }
}
