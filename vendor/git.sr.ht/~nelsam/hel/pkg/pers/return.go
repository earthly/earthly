// This is free and unencumbered software released into the public
// domain.  For more information, see <http://unlicense.org> or the
// accompanying UNLICENSE file.

package pers

import (
	"fmt"
	"reflect"

	"git.sr.ht/~nelsam/hel/sel"
)

// Return will add a given value to the channel or struct of channels.
// This isn't very useful with a single value, so it's intended more
// to support structs full of channels, such as the ones that hel
// generates for return values in its mocks.
//
// Return panics if:
// - the passed in mock value is not a valid mock field (a channel or
//   struct of channels).
// - the passed in args cannot be returned on the mock field.
// - the passed in mock value is already full and sending another
//   return value would block.
func Return(mock interface{}, args ...interface{}) {
	cases, err := sel.Cases(reflect.SelectSend, mock, args...)
	if err != nil {
		panic(fmt.Errorf("pers: returning %v on mock (%T) is not possible: %w", args, mock, err))
	}
	def := reflect.SelectCase{
		Dir: reflect.SelectDefault,
	}
	for _, c := range cases {
		defChoice := 1
		chosen, _, _ := reflect.Select([]reflect.SelectCase{c, def})
		if chosen == defChoice {
			panic(fmt.Sprintf("pers: returning %v on mock (%T) would block [hint: increase channel buffer size]", args, mock))
		}
	}
}
