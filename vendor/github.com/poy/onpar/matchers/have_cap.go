package matchers

import (
	"fmt"
	"reflect"
)

// This matcher works on Slices, Arrays, Maps and Channels and will succeed if the
// type has the specified capacity.
type HaveCapMatcher struct {
	expected int
}

// HaveCap returns a HaveCapMatcher with the specified capacity
func HaveCap(expected int) HaveCapMatcher {
	return HaveCapMatcher{
		expected: expected,
	}
}

func (m HaveCapMatcher) Match(actual any) (any, error) {
	var c int
	switch reflect.TypeOf(actual).Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		c = reflect.ValueOf(actual).Cap()
	default:
		return nil, fmt.Errorf("'%v' (%T) is not a Slice, Array, Map or Channel", actual, actual)
	}

	if c != m.expected {
		return nil, fmt.Errorf("%v (cap=%d) does not have a capacity %d", actual, c, m.expected)
	}

	return actual, nil
}
