package matchers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// MatchJSONMatcher converts both expected and actual to a map[string]any
// and does a reflect.DeepEqual between them
type MatchJSONMatcher struct {
	expected any
}

// MatchJSON returns an MatchJSONMatcher with the expected value
func MatchJSON(expected any) MatchJSONMatcher {
	return MatchJSONMatcher{
		expected: expected,
	}
}

func (m MatchJSONMatcher) Match(actual any) (any, error) {
	a, sa, err := m.unmarshal(actual)
	if err != nil {
		return nil, fmt.Errorf("Error with %s: %s", sa, err)
	}

	e, se, err := m.unmarshal(m.expected)
	if err != nil {
		return nil, fmt.Errorf("Error with %s: %s", se, err)
	}

	if !reflect.DeepEqual(a, e) {
		return nil, fmt.Errorf("expected %s to equal %s", sa, se)
	}

	return actual, nil
}

func (m MatchJSONMatcher) unmarshal(x any) (any, string, error) {
	var result any
	var s string

	switch x := x.(type) {
	case []byte:
		if err := json.Unmarshal(x, &result); err != nil {
			return nil, string(x), err
		}
		s = string(x)

	case string:
		if err := json.Unmarshal([]byte(x), &result); err != nil {
			return nil, x, err
		}
		s = x

	case *string:
		if x == nil {
			return nil, "", fmt.Errorf("*string cannot be nil")
		}
		s = *x
		if err := json.Unmarshal([]byte(s), &result); err != nil {
			return nil, s, err
		}

	default:
		return nil, "", fmt.Errorf("must be a []byte, *string, or string")
	}

	return result, s, nil
}
