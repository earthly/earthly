package matchers

type OrMatcher struct {
	Children []Matcher
}

func Or(a, b Matcher, ms ...Matcher) OrMatcher {
	return OrMatcher{
		Children: append(append([]Matcher{a}, b), ms...),
	}
}

func (m OrMatcher) Match(actual any) (any, error) {
	var err error
	for _, child := range m.Children {
		_, err = child.Match(actual)
		if err == nil {
			return actual, nil
		}
	}
	return nil, err
}
