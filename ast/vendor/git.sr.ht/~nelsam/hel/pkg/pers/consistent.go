// This is free and unencumbered software released into the public
// domain.  For more information, see <http://unlicense.org> or the
// accompanying UNLICENSE file.

package pers

import (
	"fmt"
	"reflect"
	"sync"

	"git.sr.ht/~nelsam/hel/sel"
)

// T represents the methods we need from the testing.T or testing.B types.
type T interface {
	Fatalf(string, ...interface{})
	Cleanup(func())
}

// ConsistentlyReturn will continue adding a given value to the channel until
// the test is done or the returned stop function is called, whichever happens
// first. When ConsistentlyReturn stops adding values to the channel(s), it will
// drain those channels before returning.
//
// After the first call to stop (or after the test completes), calls to stop
// will be a no-op.
//
// The value for mock may be either a channel or a struct full of channels.
// ConsistentlyReturn will panic otherwise.
//
// ConsistentlyReturn will panic if:
// - args contains a different number of arguments than the number of channels
//   on mock.
// - any of the arguments passed in are not compatible to the return types of
//   mock.
func ConsistentlyReturn(t T, mock interface{}, args ...interface{}) (stop func()) {
	cases, err := sel.Cases(reflect.SelectSend, mock, args...)
	if err != nil {
		panic(fmt.Errorf("pers: consistently returning %v on mock (%T) is not possible: %w", args, mock, err))
	}
	done := make(chan struct{})
	exited := make(chan struct{})
	go consistentlyReturn(cases, done, exited, args...)
	var once sync.Once
	stop = func() {
		once.Do(func() {
			close(done)
			<-exited
			drain(cases)
		})
	}
	t.Cleanup(stop)
	return stop
}

func drain(cases []reflect.SelectCase) {
	for _, c := range cases {
		c.Dir = reflect.SelectRecv
	}
	defIdx := len(cases)
	def := reflect.SelectCase{
		Dir: reflect.SelectDefault,
	}
	cases = append(cases, def)
	for {
		chosen, _, _ := reflect.Select(cases)
		if chosen == defIdx {
			return
		}
	}
}

func consistentlyReturn(cases []reflect.SelectCase, done, exited chan struct{}, args ...interface{}) {
	defer close(exited)
	doneIdx := len(cases)
	cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)})
	for {
		chosen, _, _ := reflect.Select(cases)
		if chosen == doneIdx {
			return
		}
	}
}
