package matchers

import "fmt"

// BeFalseMatcher will succeed if actual is false.
type BeFalseMatcher struct{}

// BeFalse will return a BeFalseMatcher
func BeFalse() BeFalseMatcher {
	return BeFalseMatcher{}
}

func (m BeFalseMatcher) Match(actual any) (any, error) {
	f, ok := actual.(bool)
	if !ok {
		return nil, fmt.Errorf("'%v' (%[1]T) is not a bool", actual)
	}

	if f {
		return nil, fmt.Errorf("%t is not false", actual)
	}

	return actual, nil
}
