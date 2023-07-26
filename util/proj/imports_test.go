package proj_test

import (
	"git.sr.ht/~nelsam/hel/pkg/pers"
	"github.com/pkg/errors"
	"github.com/poy/onpar/matchers"
)

var (
	equal        = matchers.Equal
	haveOccurred = matchers.HaveOccurred
	not          = matchers.Not

	haveMethodExecuted = pers.HaveMethodExecuted
	withArgs           = pers.WithArgs
)

type beErrMatcher struct {
	expected error
}

func beErr(err error) beErrMatcher {
	return beErrMatcher{expected: err}
}

func (m beErrMatcher) Match(actual any) (any, error) {
	err, ok := actual.(error)
	if !ok {
		return nil, errors.Errorf("expected %T to be of type error", actual)
	}
	if !errors.Is(err, m.expected) {
		return nil, errors.Errorf("expected %q to wrap error %q", err, m.expected)
	}
	return actual, nil
}
