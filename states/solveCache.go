package states

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

type SolveCache map[string]llb.State
type StateKeyFunc func() string

func NewSolveCache() *SolveCache {
	m := SolveCache(map[string]llb.State{})
	return &m
}

func (sc *SolveCache) Get(kf StateKeyFunc) (llb.State, bool) {
	s, ok := (*sc)[kf()]
	return s, ok
}

func (sc *SolveCache) Set(kf StateKeyFunc, state llb.State) {
	(*sc)[kf()] = state
}

func (sc *SolveCache) Delete(kf StateKeyFunc, state llb.State) {
	delete((*sc), kf())
}

// BuildKey builds a state key from a
func KeyFromHashAndTag(targetHash, dockerTag string) StateKeyFunc {
	return func() string {
		return fmt.Sprintf("%s-%s", dockerTag, targetHash)
	}
}

func KeyFromHash(targetHash string) StateKeyFunc {
	return func() string {
		return targetHash
	}
}
