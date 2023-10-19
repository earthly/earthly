package matchers

import "fmt"

// BeBelowMatcher accepts a float64. It succeeds if the actual is
// less than the expected.
type BeBelowMatcher struct {
	expected float64
}

// BeBelow returns a BeBelowMatcher with the expected value.
func BeBelow(expected float64) BeBelowMatcher {
	return BeBelowMatcher{
		expected: expected,
	}
}

func (m BeBelowMatcher) Match(actual any) (any, error) {
	f, err := m.toFloat(actual)
	if err != nil {
		return nil, err
	}

	if f >= m.expected {
		return nil, fmt.Errorf("%v is not below %f", actual, m.expected)
	}

	return actual, nil
}

func (m BeBelowMatcher) toFloat(actual any) (float64, error) {
	switch x := actual.(type) {
	case int:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case float32:
		return float64(x), nil
	case float64:
		return x, nil
	default:
		return 0, fmt.Errorf("Unsupported type %T", actual)
	}
}
