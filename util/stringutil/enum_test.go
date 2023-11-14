package stringutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Prefix_TestEnum int32

const (
	Prefix_TEST_ENUM_VAL1                Prefix_TestEnum = 0
	Prefix_TEST_ENUM_VAl_WITH_UNDERSCORE Prefix_TestEnum = 1
)

func (e Prefix_TestEnum) String() string {
	if e == Prefix_TEST_ENUM_VAL1 {
		return "val1"
	}
	return "val_with_underscore"
}

func Test_EnumToString(t *testing.T) {

	tests := map[string]struct {
		input    Enum
		f        EnumToStringFunc
		expected string
	}{
		"Title short value": {
			input:    Prefix_TEST_ENUM_VAL1,
			f:        Title,
			expected: "Val1",
		},
		"Title value with underscores": {
			input:    Prefix_TEST_ENUM_VAl_WITH_UNDERSCORE,
			f:        Title,
			expected: "Val With Underscore",
		},
		"Lower short value": {
			input:    Prefix_TEST_ENUM_VAL1,
			f:        Lower,
			expected: "val1",
		},
		"Lower value with underscores": {
			input:    Prefix_TEST_ENUM_VAl_WITH_UNDERSCORE,
			f:        Lower,
			expected: "val with underscore",
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
	input := []Prefix_TestEnum{Prefix_TEST_ENUM_VAL1, Prefix_TEST_ENUM_VAl_WITH_UNDERSCORE, Prefix_TEST_ENUM_VAL1}
	res := EnumToStringArray(input, Title)
	assert.Equal(t, []string{"Val1", "Val With Underscore", "Val1"}, res)
}
