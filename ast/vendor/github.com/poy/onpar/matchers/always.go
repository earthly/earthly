package matchers

import (
	"time"
)

// AlwaysMatcher matches by polling the child matcher until it returns an error.
// It will return an error the first time the child matcher returns an
// error. If the child matcher never returns an error,
// then it will return a nil.
//
// Duration is the longest scenario for the matcher
// if the child matcher continues to return nil
//
// Interval is the period between polling.
type AlwaysMatcher struct {
	Matcher            Matcher
	Duration, Interval time.Duration
}

// Always returns a default AlwaysMatcher. Length of 100ms and rate 10ms
func Always(m Matcher) AlwaysMatcher {
	return AlwaysMatcher{
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
func (m AlwaysMatcher) Match(actual any) (any, error) {
	if m.Duration == 0 {
		m.Duration = 100 * time.Millisecond
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
		if err != nil {
			return nil, err
		}

		time.Sleep(m.Interval)
	}

	return value, nil
}
