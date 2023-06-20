package flagutil

import (
	"github.com/urfave/cli/v2"
	"reflect"
	"testing"
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
