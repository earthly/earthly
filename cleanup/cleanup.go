package cleanup

// CloseFun is a cleanup function to be executed.
type CloseFun = func() error

// Collection is a collection of cleanup operations.
type Collection struct {
	closeFuns []CloseFun
}

// NewCollection returns a new Collection.
func NewCollection() *Collection {
	return &Collection{}
}

// Add adds a CloseFun to the collection.
func (c *Collection) Add(cf CloseFun) {
	c.closeFuns = append(c.closeFuns, cf)
}

// Close executes all cleanup operations.
func (c *Collection) Close() []error {
	var errs []error
	for _, cf := range c.closeFuns {
		err := cf()
		if err != nil {
			errs = append(errs, err)
			// Keep going.
		}
	}
	return errs
}
