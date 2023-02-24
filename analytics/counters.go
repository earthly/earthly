package analytics

import (
	"sync"

	"github.com/earthly/cloud-api/analytics"
)

// Counters is a threadsafe collection of counters
type Counters struct {
	counts map[string]*analytics.SendAnalyticsRequest_SubSystem
	mutex  sync.Mutex
}

var counts Counters

// Count increases the global count of (subsystem, key) which then gets reported when CollectAnalytics is called.
func (c *Counters) Count(subsystem, key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.counts == nil {
		c.counts = make(map[string]*analytics.SendAnalyticsRequest_SubSystem)
	}
	m, ok := c.counts[subsystem]
	if !ok {
		m = &analytics.SendAnalyticsRequest_SubSystem{SubSystem: make(map[string]int32)}
		c.counts[subsystem] = m
	}
	m.SubSystem[key]++
}

// getMap locks the Counter mutex, and returns the underlying counters map,
// and a function that must be called when the map is no longer being used.
func (c *Counters) getMap() (map[string]*analytics.SendAnalyticsRequest_SubSystem, func()) {
	c.mutex.Lock()
	return c.counts, c.mutex.Unlock
}
