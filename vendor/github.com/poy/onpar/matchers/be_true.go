package matchers

import "fmt"

// BeTrueMatcher will succeed if actual is true.
type BeTrueMatcher struct{}

// BeTrue will return a BeTrueMatcher
func BeTrue() BeTrueMatcher {
	return BeTrueMatcher{}
}

func (m BeTrueMatcher) Match(actual any) (any, error) {
	f, ok := actual.(bool)
	if !ok {
		return nil, fmt.Errorf("'%v' (%[1]T) is not a bool", actual)
	}

	if !f {
		return nil, fmt.Errorf("%t is not true", actual)
	}
	return actual, nil
}
