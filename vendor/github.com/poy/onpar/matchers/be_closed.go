package matchers

import (
	"fmt"
	"reflect"
)

// BeClosedMatcher only accepts a readable channel.
// It will error for anything else.
// It will succeed if the channel is closed.
type BeClosedMatcher struct{}

// BeClosed returns a BeClosedMatcher
func BeClosed() BeClosedMatcher {
	return BeClosedMatcher{}
}

func (m BeClosedMatcher) Match(actual any) (any, error) {
	t := reflect.TypeOf(actual)
	if t.Kind() != reflect.Chan || t.ChanDir() == reflect.SendDir {
		return nil, fmt.Errorf("%s is not a readable channel", t.String())
	}

	v := reflect.ValueOf(actual)

	winnerIndex, _, open := reflect.Select([]reflect.SelectCase{
		reflect.SelectCase{Dir: reflect.SelectRecv, Chan: v},
		reflect.SelectCase{Dir: reflect.SelectDefault},
	})

	if winnerIndex == 0 && !open {
		return actual, nil
	}

	return nil, fmt.Errorf("channel open")
}
