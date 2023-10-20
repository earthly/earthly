package flagutil

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationSet(t *testing.T) {

	tests := map[string]struct {
		value    string
		err      error
		expected Duration
	}{
		"parse value successfully": {
			value:    "3h",
			expected: Duration(3 * time.Hour),
		},
		"parse days value successfully": {
			value:    "5d",
			expected: Duration(5 * 24 * time.Hour),
		},
		"returns parsing error": {
			value: "5dd",
			err:   errors.New("parse error"),
		},
		"returns parsing error #2": {
			value: "1k",
			err:   errors.New("parse error"),
		},
		"empty string is 0": {},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var d Duration
			err := d.Set(tc.value)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, d)
		})

	}
}
