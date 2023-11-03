package matchers

type ChainMatcher struct {
	Children []Matcher
}

func Chain(a, b Matcher, ms ...Matcher) ChainMatcher {
	return ChainMatcher{
		Children: append(append([]Matcher{a}, b), ms...),
	}
}

func (m ChainMatcher) Match(actual any) (any, error) {
	var err error
	next := actual
	for _, child := range m.Children {
		next, err = child.Match(next)
		if err != nil {
			return nil, err
		}
	}
	return next, nil
}
