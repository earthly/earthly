package matchers

import (
	"fmt"
	"regexp"
)

type MatchRegexpMatcher struct {
	pattern string
}

func MatchRegexp(pattern string) MatchRegexpMatcher {
	return MatchRegexpMatcher{
		pattern: pattern,
	}
}

func (m MatchRegexpMatcher) Match(actual any) (any, error) {
	r, err := regexp.Compile(m.pattern)
	if err != nil {
		return nil, err
	}

	s, ok := actual.(string)
	if !ok {
		return nil, fmt.Errorf("'%v' (%T) is not a string", actual, actual)
	}

	if !r.MatchString(s) {
		return nil, fmt.Errorf("%s does not match pattern %s", s, m.pattern)
	}

	return actual, nil
}
