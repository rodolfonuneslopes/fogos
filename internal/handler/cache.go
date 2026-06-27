package handler

import (
	"sync"
	"time"

	"github.com/rodolfonuneslopes/fogos/internal/fogos"
)

const cacheTTL = 60 * time.Second

type cacheEntry struct {
	incidents []fogos.Incident
	fetchedAt time.Time
}

type incidentCache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

func newIncidentCache() *incidentCache {
	return &incidentCache{entries: make(map[string]cacheEntry)}
}

func (c *incidentCache) get(concelho string) ([]fogos.Incident, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[concelho]
	if !ok || time.Since(e.fetchedAt) > cacheTTL {
		return nil, false
	}
	return e.incidents, true
}

func (c *incidentCache) set(concelho string, incidents []fogos.Incident) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[concelho] = cacheEntry{incidents: incidents, fetchedAt: time.Now()}
}
