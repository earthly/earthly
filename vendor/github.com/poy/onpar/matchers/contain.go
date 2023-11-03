package matchers

import (
	"fmt"
	"reflect"
)

type ContainMatcher struct {
	values []any
}

func Contain(values ...any) ContainMatcher {
	return ContainMatcher{
		values: values,
	}
}

func (m ContainMatcher) Match(actual any) (any, error) {
	actualType := reflect.TypeOf(actual)
	if actualType.Kind() != reflect.Slice && actualType.Kind() != reflect.Array {
		return nil, fmt.Errorf("%s is not a Slice or Array", actualType.Kind())
	}

	actualValue := reflect.ValueOf(actual)
	for _, elem := range m.values {
		if !m.containsElem(actualValue, elem) {
			return nil, fmt.Errorf("%v does not contain %v", actual, elem)
		}
	}

	return actual, nil
}

func (m ContainMatcher) containsElem(actual reflect.Value, elem any) bool {
	for i := 0; i < actual.Len(); i++ {
		if reflect.DeepEqual(actual.Index(i).Interface(), elem) {
			return true
		}
	}

	return false
}
