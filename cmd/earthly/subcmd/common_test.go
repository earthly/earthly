package subcmd

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"testing"
	"time"
)

func TestDurationSet(t *testing.T) {

	tests := map[string]struct {
		value    string
		err      error
		expected time.Duration
	}{
		"parse value successfully": {
			value:    "3h",
			expected: 3 * time.Hour,
		},
		"parse days value successfully": {
			value:    "5d",
			expected: 5 * 24 * time.Hour,
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
			d := &duration{}
			err := d.Set(tc.value)
			if tc.err != nil {
				assert.EqualError(t, err, tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, d.Value())
		})

	}
}

func TestDurationValue(t *testing.T) {
	d := &duration{}
	err := d.Set("13m")
	require.NoError(t, err)
	assert.Equal(t, 13*time.Minute, d.Value())
}
