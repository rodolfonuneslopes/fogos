package handler

import (
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const cacheTTL = 60 * time.Second

type cacheEntry struct {
	data      []byte // pre-serialized JSON
	fetchedAt time.Time
}

type incidentCache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
	group   singleflight.Group
}

func newIncidentCache() *incidentCache {
	return &incidentCache{entries: make(map[string]cacheEntry)}
}

// get returns the cached JSON bytes if the entry exists and is still fresh.
func (c *incidentCache) get(concelho string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[concelho]
	if !ok || time.Since(e.fetchedAt) > cacheTTL {
		return nil, false
	}
	return e.data, true
}

func (c *incidentCache) set(concelho string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[concelho] = cacheEntry{data: data, fetchedAt: time.Now()}
}
