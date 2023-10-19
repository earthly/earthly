package matchers

import (
	"fmt"
	"reflect"
	"time"
)

// ViaPollingMatcher matches by polling the child matcher until
// it returns a success. It will return success the first time
// the child matcher returns a success. If the child matcher
// never returns a nil, then it will return the last error.
//
// Duration is the worst case scenario for the matcher
// if the child matcher continues to return an error
//
// Interval is the period between polling.
type ViaPollingMatcher struct {
	Matcher            Matcher
	Duration, Interval time.Duration
}

// ViaPolling returns the default ViaPollingMatcher. Length of 1s
// and Rate of 10ms
func ViaPolling(m Matcher) ViaPollingMatcher {
	return ViaPollingMatcher{
		Matcher: m,
	}
}

// Match takes a value that can change over time. Therefore, the only
// two valid options are a function with no arguments and a single return
// type, or a readable channel. Anything else will return an error.
//
// If actual is a channel, then the child matcher will have to handle
// reading from the channel.
//
// If the actual is a function, then the matcher will invoke the value
// and pass the returned value to the child matcher.
func (m ViaPollingMatcher) Match(actual any) (any, error) {
	if m.Duration == 0 {
		m.Duration = time.Second
	}

	if m.Interval == 0 {
		m.Interval = 10 * time.Millisecond
	}

	f, err := fetchFunc(actual)
	if err != nil {
		return nil, err
	}

	var value any
	for i := 0; i < int(m.Duration/m.Interval); i++ {
		value, err = m.Matcher.Match(f())
		if err == nil {
			return value, nil
		}

		time.Sleep(m.Interval)
	}

	return nil, err
}

func fetchFunc(actual any) (func() any, error) {
	t := reflect.TypeOf(actual)
	switch t.Kind() {
	case reflect.Func:
		return fetchFuncFromFunc(actual)
	case reflect.Chan:
		return fetchFuncFromChan(actual)
	default:
		return nil, fmt.Errorf("invalid type: %v", t)
	}
}

func fetchFuncFromChan(actual any) (func() any, error) {
	t := reflect.TypeOf(actual)
	if t.ChanDir() == reflect.SendDir {
		return nil, fmt.Errorf("channel must be able to receive")
	}

	return func() any {
		return actual
	}, nil
}

func fetchFuncFromFunc(actual any) (func() any, error) {
	t := reflect.TypeOf(actual)
	if t.NumIn() != 0 {
		return nil, fmt.Errorf("func must not take any arguments")
	}

	if t.NumOut() != 1 {
		return nil, fmt.Errorf("func must have one return type")
	}

	return func() any {
		v := reflect.ValueOf(actual)
		retValues := v.Call(nil)
		return retValues[0].Interface()
	}, nil
}
