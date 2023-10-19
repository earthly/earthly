package matchers

import "errors"

// PanicMatcher accepts a function. It succeeds if the function panics.
type PanicMatcher struct {
}

// Panic returns a Panic matcher.
func Panic() PanicMatcher {
	return PanicMatcher{}
}

func (m PanicMatcher) Match(actual any) (result any, err error) {
	f, ok := actual.(func())
	if !ok {
		return nil, errors.New("actual must be a func()")
	}

	defer func() {
		r := recover()
		if r == nil {
			err = errors.New("expected to panic")
		}
	}()

	f()

	return nil, nil
}
