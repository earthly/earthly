package subcmd

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDurationSet(t *testing.T) {

	tests := map[string]struct {
		value    string
		err      error
		expected duration
	}{
		"parse value successfully": {
			value:    "3h",
			expected: duration(3 * time.Hour),
		},
		"parse days value successfully": {
			value:    "5d",
			expected: duration(5 * 24 * time.Hour),
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
			var d duration
			err := d.Set(tc.value)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, d)
		})

	}
}
