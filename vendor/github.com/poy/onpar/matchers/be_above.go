package matchers

import "fmt"

// BeAboveMatcher accepts a float64. It succeeds if the
// actual is greater than the expected.
type BeAboveMatcher struct {
	expected float64
}

// BeAbove returns a BeAboveMatcher with the expected value.
func BeAbove(expected float64) BeAboveMatcher {
	return BeAboveMatcher{
		expected: expected,
	}
}

func (m BeAboveMatcher) Match(actual any) (any, error) {
	f, err := m.toFloat(actual)
	if err != nil {
		return nil, err
	}

	if f <= m.expected {
		return nil, fmt.Errorf("%v is not above %f", actual, m.expected)
	}

	return actual, nil
}

func (m BeAboveMatcher) toFloat(actual any) (float64, error) {
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
