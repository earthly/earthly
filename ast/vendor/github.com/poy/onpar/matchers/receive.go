package matchers

import (
	"fmt"
	"reflect"
	"time"
)

// ReceiveOpt is an option that can be passed to the
// ReceiveMatcher constructor.
type ReceiveOpt func(ReceiveMatcher) ReceiveMatcher

// ReceiveWait is an option that makes the ReceiveMatcher
// wait for values for the provided duration before
// deciding that the channel failed to receive.
func ReceiveWait(t time.Duration) ReceiveOpt {
	return func(m ReceiveMatcher) ReceiveMatcher {
		m.timeout = t
		return m
	}
}

// ReceiveMatcher only accepts a readable channel. It will error for anything else.
// It will attempt to receive from the channel but will not block.
// It fails if the channel is closed.
type ReceiveMatcher struct {
	timeout time.Duration
}

// Receive will return a ReceiveMatcher
func Receive(opts ...ReceiveOpt) ReceiveMatcher {
	var m ReceiveMatcher
	for _, opt := range opts {
		m = opt(m)
	}
	return m
}

func (m ReceiveMatcher) Match(actual any) (any, error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Chan || t.ChanDir() == reflect.SendDir {
		return nil, fmt.Errorf("%s is not a readable channel", t.String())
	}

	timeout := reflect.SelectCase{
		Dir: reflect.SelectDefault,
	}
	if m.timeout != 0 {
		timeout.Dir = reflect.SelectRecv
		timeout.Chan = reflect.ValueOf(time.After(m.timeout))
	}
	i, v, ok := reflect.Select([]reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(actual)},
		timeout,
	})
	if i == 1 || !ok {
		return nil, fmt.Errorf("did not receive")
	}

	return v.Interface(), nil
}
