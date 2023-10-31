package flagutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildArgMatrix(t *testing.T) {
	var tests = []struct {
		in  []string
		out [][]string
	}{
		{[]string{}, [][]string{nil}},
		{[]string{"a=1"}, [][]string{{"a=1"}}},
		{[]string{"a=1", "a=2", "a=3"}, [][]string{{"a=1"}, {"a=2"}, {"a=3"}}},
		{[]string{"a=1", "b=2"}, [][]string{{"a=1", "b=2"}}},
		{[]string{"a=1", "a=3", "b=2"}, [][]string{{"a=1", "b=2"}, {"a=3", "b=2"}}},
		{[]string{"a=1", "a=3", "b=2", "b=4"}, [][]string{{"a=1", "b=2"}, {"a=1", "b=4"}, {"a=3", "b=2"}, {"a=3", "b=4"}}},
		{[]string{"a=1", "b=2", "a=3", "b=4"}, [][]string{{"a=1", "b=2"}, {"a=1", "b=4"}, {"a=3", "b=2"}, {"a=3", "b=4"}}},
		{[]string{"a=1", "b=2", "a=3", "b=4", "c=10"}, [][]string{{"a=1", "b=2", "c=10"}, {"a=1", "b=4", "c=10"}, {"a=3", "b=2", "c=10"}, {"a=3", "b=4", "c=10"}}},
		{[]string{"a=1", "a=3", "a=7", "c=10"}, [][]string{{"a=1", "c=10"}, {"a=3", "c=10"}, {"a=7", "c=10"}}},
	}

	for _, tt := range tests {
		ans, err := BuildArgMatrix(tt.in)
		assert.NoError(t, err)
		assert.Equal(t, tt.out, ans)
	}
}
