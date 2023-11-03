package matchers

import (
	"fmt"
	"reflect"
)

// HaveKeyMatcher accepts map types and will succeed if the map contains the
// specified key.
type HaveKeyMatcher struct {
	key any
}

// HaveKey returns a HaveKeyMatcher with the specified key.
func HaveKey(key any) HaveKeyMatcher {
	return HaveKeyMatcher{
		key: key,
	}
}

func (m HaveKeyMatcher) Match(actual any) (any, error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Map {
		return nil, fmt.Errorf("'%v' (%T) is not a Map", actual, actual)
	}

	if t.Key() != reflect.TypeOf(m.key) {
		return nil, fmt.Errorf("'%v' (%T) has a Key type of %v not %T", actual, actual, t.Key(), m.key)
	}

	v := reflect.ValueOf(actual)
	value := v.MapIndex(reflect.ValueOf(m.key))
	if !value.IsValid() {
		return nil, fmt.Errorf("unable to find key %v in %v", m.key, actual)
	}

	return value.Interface(), nil
}
