package llbfactory

import (
	"github.com/earthly/earthly/util/llbutil/pllb"

	"github.com/moby/buildkit/client/llb"
)

// Factory is used for constructing llb states
type Factory interface {
	// Construct creates a pllb.State
	Construct() pllb.State
}

// PreconstructedFactory holds a preconstructed pllb.State for cases
// where a factory is overkill.
type PreconstructedFactory struct {
	preconstructedState pllb.State
}

// LocalFactory holds data which can be used to create a pllb.Local state
type LocalFactory struct {
	name          string
	sharedKeyHint string
	opts          []llb.LocalOption
}

// PreconstructedState returns a pseudo-factory which returns
// the passed in state when Construct() is called.
// It is provided for cases where a factory is overkill.
func PreconstructedState(state pllb.State) Factory {
	return &PreconstructedFactory{
		preconstructedState: state,
	}
}

// Construct returns the preconstructed state that was passed to PreconstructedState()
func (f *PreconstructedFactory) Construct() pllb.State {
	return f.preconstructedState
}

// Local eventually creates a llb.Local
func Local(name string, opts ...llb.LocalOption) Factory {
	return &LocalFactory{
		name: name,
		opts: opts,
	}
}

// Copy makes a new copy of the localFactory
func (f *LocalFactory) Copy() *LocalFactory {
	newOpts := append([]llb.LocalOption{}, f.opts...)
	return &LocalFactory{
		name: f.name,
		opts: newOpts,
	}
}

// GetName returns the name of the pllb.Local state that will
// eventually be created
func (f *LocalFactory) GetName() string {
	return f.name
}

// GetSharedKey returns the set shared cache key of the pllb.Local state that will
// eventually be created
func (f *LocalFactory) GetSharedKey() string {
	return f.sharedKeyHint
}

// WithInclude adds include patterns to the factory's llb options
func (f *LocalFactory) WithInclude(patterns []string) *LocalFactory {
	f = f.Copy()
	f.opts = append(f.opts, llb.IncludePatterns(patterns))
	return f
}

// WithSharedKeyHint adds a shared key hint to the factory's llb options
func (f *LocalFactory) WithSharedKeyHint(key string) *LocalFactory {
	f = f.Copy()
	f.opts = append(f.opts, llb.SharedKeyHint(key))
	f.sharedKeyHint = key
	return f
}

// Construct constructs the pllb.Local state
func (f *LocalFactory) Construct() pllb.State {
	return pllb.Local(f.name, f.opts...)
}
