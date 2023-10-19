package matchers

import (
	"fmt"
	"strings"
)

// ContainSubstringMatcher accepts a string and succeeds
// if the actual string contains the expected string.
type ContainSubstringMatcher struct {
	substr string
}

// ContainSubstring returns a ContainSubstringMatcher with the
// expected substring.
func ContainSubstring(substr string) ContainSubstringMatcher {
	return ContainSubstringMatcher{
		substr: substr,
	}
}

func (m ContainSubstringMatcher) Match(actual any) (any, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.Contains(s, m.substr) {
		return nil, fmt.Errorf("%s does not contain %s", s, m.substr)
	}

	return actual, nil
}
