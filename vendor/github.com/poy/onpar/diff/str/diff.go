package str

type Type int

const (
	TypeMatch = iota
	TypeReplace
)

// DiffSection represents a single chunk of a diff.
type DiffSection struct {
	Type             Type
	Actual, Expected []rune
}

// Diff represents a full diff of two values.
type Diff interface {
	// Cost is the calculated cost of changing from one value to another.
	// Basically, if provided with multiple diffs, the Differ will always prefer
	// the lowest cost.
	//
	// Generally, a cost of 0 should represent exactly equal values, so negative
	// numbers shouldn't usually be used. However, if they are used, they will
	// work the same as positive values, being preferred over any value higher
	// than them.
	Cost() float64

	// Sections returns all of the sections of the diff. This will be used to
	// generate output, depending on the diff formats being used.
	Sections() []DiffSection
}
