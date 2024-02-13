package stringutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessParamsAndQuotes(t *testing.T) {
	var tests = []struct {
		in   []string
		args []string
	}{
		{[]string{}, []string{}},
		{[]string{""}, []string{""}},
		{[]string{"abc", "def", "ghi"}, []string{"abc", "def", "ghi"}},
		{[]string{"hello ", "wor(", "ld)"}, []string{"hello ", "wor( ld)"}},
		{[]string{"hello ", "(wor(", "ld)"}, []string{"hello ", "(wor( ld)"}},
		{[]string{"hello ", "\"(wor(\"", "ld)"}, []string{"hello ", "\"(wor(\"", "ld)"}},
		{[]string{"let's", "go"}, []string{"let's go"}},
		{[]string{"(hello)"}, []string{"(hello)"}},
		{[]string{"  (hello)"}, []string{"  (hello)"}},
		{[]string{"(hello", "    ooo)"}, []string{"(hello     ooo)"}},
		{[]string{"--load=(+a-test-image", "--name=foo", "--var", "bar)"}, []string{"--load=(+a-test-image --name=foo --var bar)"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(strings.Join(tt.in, " "), func(t *testing.T) {
			t.Parallel()
			actualArgs := ProcessParamsAndQuotes(tt.in)
			assert.Equal(t, tt.args, actualArgs)
		})

	}
}
