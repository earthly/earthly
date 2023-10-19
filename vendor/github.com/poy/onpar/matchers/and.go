package matchers

type AndMatcher struct {
	Children []Matcher
}

func And(a, b Matcher, ms ...Matcher) AndMatcher {
	return AndMatcher{
		Children: append(append([]Matcher{a}, b), ms...),
	}
}

func (m AndMatcher) Match(actual any) (any, error) {
	var err error
	for _, child := range m.Children {
		_, err = child.Match(actual)
		if err != nil {
			return nil, err
		}
	}
	return actual, nil
}
