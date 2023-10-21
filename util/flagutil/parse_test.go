package flagutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

func TestSplitFlagString(t *testing.T) {
	type args struct {
		value cli.StringSlice
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "passing flag multiple times",
			args: args{
				value: *(cli.NewStringSlice("a b")),
			},
			want: []string{"a", "b"},
		},
		{
			name: "passing values with a comma",
			args: args{
				value: *(cli.NewStringSlice("a,b")),
			},
			want: []string{"a", "b"},
		},
		{
			name: "passing values with a comma and multiple flags",
			args: args{
				value: *(cli.NewStringSlice("a b,c   d")),
			},
			want: []string{"a", "b", "c", "d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitFlagString(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitFlagString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseParams(t *testing.T) {
	var tests = []struct {
		in    string
		first string
		args  []string
	}{
		{"(+target/art --flag=something)", "+target/art", []string{"--flag=something"}},
		{"(+target/art --flag=something\"\")", "+target/art", []string{"--flag=something\"\""}},
		{"( \n  +target/art \t \n --flag=something\t   )", "+target/art", []string{"--flag=something"}},
		{"(+target/art --flag=something\\ --another=something)", "+target/art", []string{"--flag=something\\ --another=something"}},
		{"(+target/art --flag=something --another=something)", "+target/art", []string{"--flag=something", "--another=something"}},
		{"(+target/art --flag=\"something in quotes\")", "+target/art", []string{"--flag=\"something in quotes\""}},
		{"(+target/art --flag=\\\"something --not=in-quotes\\\")", "+target/art", []string{"--flag=\\\"something", "--not=in-quotes\\\""}},
		{"(+target/art --flag=look-ma-a-\\))", "+target/art", []string{"--flag=look-ma-a-\\)"}},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			actualFirst, actualArgs, err := ParseParams(tt.in)
			assert.NoError(t, err)
			assert.Equal(t, tt.first, actualFirst)
			assert.Equal(t, tt.args, actualArgs)
		})

	}
}

func TestNegativeParseParams(t *testing.T) {
	var tests = []struct {
		in string
	}{
		{"+target/art --flag=something)"},
		{"(+target/art --flag=something"},
		{"(+target/art --flag=\"something)"},
		{"(+target/art --flag=something\\)"},
		{"()"},
		{"(          \t\n   )"},
	}

	for _, tt := range tests {
		_, _, err := ParseParams(tt.in)
		assert.Error(t, err)
	}
}
