package states

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

// SolveCache is a formal version of the cache we keep mapping targets to their LLB state.
type SolveCache map[StateKey]llb.State

// StateKey is a type for a key in SolveCache. These keys seem to be highly convention based,
// and used elsewhere too (LocalFolders?). so this is a step atformalizing that convention,
// since we sometimes need one key, and sometimes another. It may give us some toeholds to
// help with some refactoring later.
type StateKey string

// NewSolveCache gives a new SolveCachemap instance
func NewSolveCache() *SolveCache {
	m := SolveCache(map[StateKey]llb.State{})
	return &m
}

// Get gets a LLB state out of a given solve cache, using the KeyFunc to derive the key
func (sc *SolveCache) Get(sk StateKey) (llb.State, bool) {
	s, ok := (*sc)[sk]
	return s, ok
}

// Set puts a LLB state in a given solve cache, using the KeyFunc to derive the key
func (sc *SolveCache) Set(sk StateKey, state llb.State) {
	(*sc)[sk] = state
}

// Delete removes a LLB state from a given solve cache, using the KeyFunc to derive the key
func (sc *SolveCache) Delete(sk StateKey, state llb.State) {
	delete((*sc), sk)
}

// KeyFromHashAndTag builds a state key from a given target state and a docker tag.
// This is useful when you want to reference the same image but with a different name.
func KeyFromHashAndTag(target *SingleTarget, dockerTag string) (StateKey, error) {
	hash, err := target.TargetInput.Hash()
	if err != nil {
		return StateKey(""), errors.Wrap(err, "target input hash")
	}

	key := fmt.Sprintf("%s-%s", dockerTag, hash)
	return StateKey(key), nil
}

// KeyFromState is a simple wrapper to get a key from a given state using the hash of its target.
func KeyFromState(target *SingleTarget) (StateKey, error) {
	hash, err := target.TargetInput.Hash()
	if err != nil {
		return StateKey(""), errors.Wrap(err, "target input hash")
	}

	return StateKey(hash), nil
}
