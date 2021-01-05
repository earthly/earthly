package states

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

// SolveCache is a formal version of the cache we keep mapping targets to their LLB state.
type SolveCache map[string]llb.State

// StateKeyFunc is a function that can build a key for the SolveCache. These keys seem to be
// Highly convention based, and used elsewhere too (LocalFolders?), so this is a step at
// formalizing that convention, since we sometimes need one key, and sometimes another.
// It may give us some toeholds to help with some refactoring later.
type StateKeyFunc func() string

// NewSolveCache gives a new SolveCachemap instance
func NewSolveCache() *SolveCache {
	m := SolveCache(map[string]llb.State{})
	return &m
}

// Get gets a LLB state out of a given solve cache, using the KeyFunc to derive the key
func (sc *SolveCache) Get(kf StateKeyFunc) (llb.State, bool) {
	s, ok := (*sc)[kf()]
	return s, ok
}

// Set puts a LLB state in a given solve cache, using the KeyFunc to derive the key
func (sc *SolveCache) Set(kf StateKeyFunc, state llb.State) {
	(*sc)[kf()] = state
}

// Delete removes a LLB state from a given solve cache, using the KeyFunc to derive the key
func (sc *SolveCache) Delete(kf StateKeyFunc, state llb.State) {
	delete((*sc), kf())
}

// KeyFromHashAndTag builds a state key from a given target state and a docker tag.
// This is useful when you want to reference the same image but with a different name.
func KeyFromHashAndTag(target *SingleTarget, dockerTag string) (StateKeyFunc, error) {
	// Keys must remain stable... kinda the point here
	hash, err := target.TargetInput.Hash()
	if err != nil {
		return nil, errors.Wrap(err, "target input hash")
	}

	return func() string {
		return fmt.Sprintf("%s-%s", dockerTag, hash)
	}, nil
}

// KeyFromHash is a simple wrapper to get a key from a given state.
func KeyFromHash(targetHash string) StateKeyFunc {
	return func() string {
		return targetHash
	}
}
