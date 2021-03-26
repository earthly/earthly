package states

import (
	"context"
	"fmt"

	"github.com/earthly/earthly/syncutil/synccache"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

// SolveCacheConstructor is func taking a StateKey and returning a state.
type SolveCacheConstructor func(context.Context, StateKey) (llb.State, error)

// SolveCache is a formal version of the cache we keep mapping targets to their LLB state.
type SolveCache struct {
	store *synccache.SyncCache // StateKey -> llb.State
}

// StateKey is a type for a key in SolveCache. These keys seem to be highly convention based,
// and used elsewhere too (LocalFolders?). so this is a step at formalizing that convention,
// since we sometimes need one key, and sometimes another. It may give us some toeholds to
// help with some refactoring later.
type StateKey string

// NewSolveCache gives a new SolveCachemap instance
func NewSolveCache() *SolveCache {
	return &SolveCache{
		store: synccache.New(),
	}
}

// Do sets an LLB state in the given solve cache. If the state has been previously constructed,
// it is returned immediately without calling the constructor again.
func (sc *SolveCache) Do(ctx context.Context, sk StateKey, constructor SolveCacheConstructor) (llb.State, error) {
	stateValue, err := sc.store.Do(ctx, sk, func(ctx context.Context, k interface{}) (interface{}, error) {
		return constructor(ctx, k.(StateKey))
	})
	if err != nil {
		return llb.State{}, err
	}
	return stateValue.(llb.State), nil
}

// KeyFromHashAndTag builds a state key from a given target state and a docker tag.
// This is useful when you want to reference the same image but with a different name.
func KeyFromHashAndTag(target *SingleTarget, dockerTag string) (StateKey, error) {
	hash, err := target.TargetInput().Hash()
	if err != nil {
		return StateKey(""), errors.Wrap(err, "target input hash")
	}

	key := fmt.Sprintf("%s-%s", dockerTag, hash)
	return StateKey(key), nil
}

// KeyFromState is a simple wrapper to get a key from a given state using the hash of its target.
func KeyFromState(target *SingleTarget) (StateKey, error) {
	hash, err := target.TargetInput().Hash()
	if err != nil {
		return StateKey(""), errors.Wrap(err, "target input hash")
	}

	return StateKey(hash), nil
}
