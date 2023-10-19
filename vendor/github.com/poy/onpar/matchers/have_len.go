package matchers

import (
	"fmt"
	"reflect"
)

// HaveLenMatcher accepts Strings, Slices, Arrays, Maps and Channels. It will
// succeed if the type has the specified length.
type HaveLenMatcher struct {
	expected int
}

// HaveLen returns a HaveLenMatcher with the specified length.
func HaveLen(expected int) HaveLenMatcher {
	return HaveLenMatcher{
		expected: expected,
	}
}

func (m HaveLenMatcher) Match(actual any) (any, error) {
	var l int
	switch reflect.TypeOf(actual).Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String, reflect.Chan:
		l = reflect.ValueOf(actual).Len()
	default:
		return nil, fmt.Errorf("'%v' (%T) is not a Slice, Array, Map, String or Channel", actual, actual)
	}

	if l != m.expected {
		return nil, fmt.Errorf("%v (len=%d) does not have a length of %d", actual, l, m.expected)
	}

	return actual, nil
}
