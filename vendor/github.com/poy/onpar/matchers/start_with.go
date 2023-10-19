package matchers

import (
	"fmt"
	"strings"
)

// StartWithMatcher accepts a string and succeeds
// if the actual string starts with the expected string.
type StartWithMatcher struct {
	prefix string
}

// StartWith returns a StartWithMatcher with the expected prefix.
func StartWith(prefix string) StartWithMatcher {
	return StartWithMatcher{
		prefix: prefix,
	}
}

func (m StartWithMatcher) Match(actual any) (any, error) {
	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !strings.HasPrefix(s, m.prefix) {
		return nil, fmt.Errorf("%s does not start with %s", s, m.prefix)
	}

	return actual, nil
}
