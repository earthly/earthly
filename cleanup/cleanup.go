package cleanup

import "sync"

// CloseFun is a cleanup function to be executed.
type CloseFun = func() error

// Collection is a collection of cleanup operations.
type Collection struct {
	mu        sync.Mutex
	closeFuns []CloseFun
}

// NewCollection returns a new Collection.
func NewCollection() *Collection {
	return &Collection{}
}

// Add adds a CloseFun to the collection.
func (c *Collection) Add(cf CloseFun) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closeFuns = append(c.closeFuns, cf)
}

// Close executes all cleanup operations and empties the collection.
func (c *Collection) Close() []error {
	c.mu.Lock()
	cfs := c.closeFuns
	c.closeFuns = []CloseFun{}
	c.mu.Unlock()
	var errs []error
	for _, cf := range cfs {
		err := cf()
		if err != nil {
			errs = append(errs, err)
			// Keep going.
		}
	}
	return errs
}
