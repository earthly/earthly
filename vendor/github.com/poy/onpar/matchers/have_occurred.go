package matchers

import "fmt"

// HaveOccurredMatcher will succeed if the actual value is a non-nil error.
type HaveOccurredMatcher struct {
}

// HaveOccurred returns a HaveOccurredMatcher
func HaveOccurred() HaveOccurredMatcher {
	return HaveOccurredMatcher{}
}

func (m HaveOccurredMatcher) Match(actual any) (any, error) {
	e, ok := actual.(error)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not an error", actual, actual)
	}

	if e == nil {
		return nil, fmt.Errorf("err to not be nil")
	}

	return nil, nil
}
