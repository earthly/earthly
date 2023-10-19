package vegr

import (
	"fmt"
	"reflect"
	"time"

	"git.sr.ht/~nelsam/hel/sel"
)

type boolSet []bool

func newBoolSet(length int) boolSet {
	return make(boolSet, length)
}

func (s boolSet) all() bool {
	for _, v := range s {
		if !v {
			return false
		}
	}
	return true
}

type T interface {
	Helper()
	Failed() bool
	Fatalf(string, ...interface{})
}

func PopulateReturns(t T, name string, timeout time.Duration, mock interface{}, addrs ...interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			if t.Failed() {
				return
			}
			panic(r)
		}
	}()
	cases, err := sel.Cases(reflect.SelectRecv, mock)
	if err != nil {
		panic(fmt.Errorf("hel: PopulateReturns was called with incorrect mock type (%T): %w", mock, err))
	}
	if len(addrs) != len(cases) {
		panic(fmt.Errorf("hel: PopulateReturns was called with %d channels but only %d addresses", len(cases), len(addrs)))
	}
	var vals []reflect.Value
	for _, a := range addrs {
		v := reflect.ValueOf(a)
		if v.Kind() != reflect.Ptr {
			panic(fmt.Errorf("hel: PopulateReturns was called with non-pointer type (%T)", a))
		}
		if v.IsNil() {
			panic(fmt.Errorf("hel: PopulateReturns was called with nil pointer of type (%T)", a))
		}
		vals = append(vals, v.Elem())
	}
	done := newBoolSet(len(cases))
	deadline := time.NewTimer(timeout)
	defer deadline.Stop()

	timeoutIdx := len(cases)
	cases = append(cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(deadline.C),
	})
	if t.Failed() {
		return
	}
	for !done.all() {
		chosen, recv, ok := reflect.Select(cases)
		if !ok {
			panic(fmt.Errorf("hel: PopulateReturns called on closed mock (type %T)", mock))
		}
		if t.Failed() {
			return
		}
		if chosen == timeoutIdx {
			defer func() {
				if r := recover(); r != nil {
					panic(fmt.Errorf("hel: panic calling t.Fatalf from mock method %v: %v", name, r))
				}
			}()
			t.Fatalf("hel: mock method %v timed out after %v waiting for return on type (%T)", name, timeout, mock)
			return
		}
		vals[chosen].Set(recv)
		done[chosen] = true
	}
}
