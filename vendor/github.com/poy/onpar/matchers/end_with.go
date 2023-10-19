package matchers

import (
	"fmt"
	"strings"
)

// EndWithMatcher accepts a string and succeeds
// if the actual string ends with the expected string.
type EndWithMatcher struct {
	suffix string
}

// EndWith returns an EndWithMatcher with the expected suffix.
func EndWith(suffix string) EndWithMatcher {
	return EndWithMatcher{
		suffix: suffix,
	}
}

func (m EndWithMatcher) Match(actual any) (any, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.HasSuffix(s, m.suffix) {
		return nil, fmt.Errorf("%s does not end with %s", s, m.suffix)
	}

	return actual, nil
}
