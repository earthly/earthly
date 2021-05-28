package analytics

import (
	"sync"
)

// Counters is a threadsafe collection of counters
type Counters struct {
	counts map[string]map[string]int
	mutex  sync.Mutex
}

var counts Counters

// Count increases the global count of (subsystem, key) which then gets reported when CollectAnalytics is called.
func (c *Counters) Count(subsystem, key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.counts == nil {
		c.counts = map[string]map[string]int{}
	}
	m, ok := c.counts[subsystem]
	if !ok {
		m = map[string]int{}
		c.counts[subsystem] = m
	}
	m[key]++
}

// getMap locks the Counter mutex, and returns the underlying counters map,
// and a function that must be called when the map is no longer being used.
func (c *Counters) getMap() (map[string]map[string]int, func()) {
	c.mutex.Lock()
	return c.counts, c.mutex.Unlock
}
