package stringutil

import (
	"testing"

	"github.com/earthly/cloud-api/logstream"
	"github.com/stretchr/testify/assert"
)

func Test_EnumToString(t *testing.T) {

	tests := map[string]struct {
		input    ProtoEnum
		f        EnumToStringFunc
		expected string
	}{
		"Title": {
			input:    logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED,
			f:        Title,
			expected: "Buildkit Crashed",
		},
		"Lower": {
			input:    logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED,
			f:        Lower,
			expected: "buildkit crashed",
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			res := tc.f(tc.input)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func Test_EnumToStringArray(t *testing.T) {
	input := []ProtoEnum{logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED, logstream.FailureType_FAILURE_TYPE_CONNECTION_FAILURE}
	res := EnumToStringArray(input, Title)
	assert.Equal(t, []string{"Buildkit Crashed", "Connection Failure"}, res)
}
