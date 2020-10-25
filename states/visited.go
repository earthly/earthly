package states

// VisitedCollection is a collection of visited targets.
type VisitedCollection struct {
	Visited map[string][]*SingleTarget
	// Same collection as above, but as a list, to make the ordering consistent.
	VisitedList []*SingleTarget
}

// NewVisitedCollection returns a collection of visited targets.
func NewVisitedCollection() *VisitedCollection {
	return &VisitedCollection{
		Visited: make(map[string][]*SingleTarget),
	}
}

// Add adds a target to the collection.
func (vc *VisitedCollection) Add(target string, sts *SingleTarget) {
	vc.Visited[target] = append(vc.Visited[target], sts)
	vc.VisitedList = append(vc.VisitedList, sts)
}
