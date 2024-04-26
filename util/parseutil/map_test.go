package parseutil

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToMap(t *testing.T) {
	tests := map[string]struct {
		input       string
		expected    map[string]string
		expectedErr error
	}{
		"happy path": {
			input: "key1=val1,key2= val2 , key3 =val3 ,",
			expected: map[string]string{
				"key1": "val1",
				"key2": "val2",
				"key3": "val3",
			},
		},
		"happy path - empty map": {
			input:    "   ",
			expected: map[string]string{},
		},
		"happy path - single value": {
			input: "key1=val1",
			expected: map[string]string{
				"key1": "val1",
			},
		},
		"error when no equal sign": {
			input:       "key1=val1,key2 val2 , key3 =val3 ,",
			expectedErr: errors.New("key/value must be set with ="),
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res, err := StringToMap(tc.input)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, res)
		})
	}
}
