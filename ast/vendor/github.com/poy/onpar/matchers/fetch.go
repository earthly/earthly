package matchers

import (
	"fmt"
	"reflect"
)

type FetchMatcher struct {
	OutputTo any
}

func Fetch(outputTo any) FetchMatcher {
	return FetchMatcher{
		OutputTo: outputTo,
	}
}

func (m FetchMatcher) Match(actual any) (any, error) {
	outType := reflect.TypeOf(m.OutputTo)
	outValue := reflect.ValueOf(m.OutputTo)
	actualValue := reflect.ValueOf(actual)

	if outType.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("%s is not a pointer type", outType.String())
	}

	if !reflect.TypeOf(actualValue.Interface()).AssignableTo(outType.Elem()) {
		return nil, fmt.Errorf("can not assigned %s to %s",
			reflect.TypeOf(actualValue.Interface()).String(),
			reflect.TypeOf(m.OutputTo).String(),
		)
	}

	outValue.Elem().Set(actualValue)
	return actual, nil
}
