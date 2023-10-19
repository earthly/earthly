package matchers

import (
	"fmt"
	"reflect"
)

// IsNilMatcher will succeed if actual is nil.
type IsNilMatcher struct{}

// IsNil will return a IsNilMatcher.
func IsNil() IsNilMatcher {
	return IsNilMatcher{}
}

func (m IsNilMatcher) Match(actual any) (any, error) {
	if actual == nil {
		return nil, nil
	}

	var isNil bool
	switch reflect.TypeOf(actual).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		isNil = reflect.ValueOf(actual).IsNil()
	}

	if isNil {
		return nil, nil
	}

	return actual, fmt.Errorf("%v is not nil", actual)
}
