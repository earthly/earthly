package inputgraph

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_containsShellExpr(t *testing.T) {
	cases := []struct {
		desc string
		val  string
		want bool
	}{
		{
			desc: "single",
			val:  "$(echo 'hello')",
			want: true,
		},
		{
			desc: "nested",
			val:  "$(echo -n $(cat /tmp/x))",
			want: true,
		},
		{
			desc: "invalid 1",
			val:  "$($()",
			want: false,
		},
		{
			desc: "invalid 2",
			val:  ")$(",
			want: false,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			got := containsShellExpr(c.val)
			require.Equal(t, c.want, got)
		})
	}
}

func Test_evalConditions(t *testing.T) {
	cases := []struct {
		val  string
		want [2]bool
	}{
		{
			val:  `[ -f "file" ]`,
			want: [2]bool{false, false},
		},
		{
			val:  "[ true ]",
			want: [2]bool{true, true},
		},
		{
			val:  "[ false ]",
			want: [2]bool{false, true},
		},
		{
			val:  `[ -n "foo" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ -n "" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ -z "" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ -z "foo" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ "foo" == "bar" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ "foo" = "bar" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ "foo" != "bar" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ "foo" > "bar" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ "foo" < "bar" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ "2" -eq "2" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ "2" -eq "4" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ "2" -ne "4" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ "2" -gt "1" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ "2" -gt "10" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ "2" -gt "foo" ]`,
			want: [2]bool{false, false},
		},
		{
			val:  `[ "2" -gt "foo" ]`,
			want: [2]bool{false, false},
		},
		{
			val:  `[ "2" -lt "100" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ "2" -lt "1" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ "2" -le "2" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ "4" -le "2" ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ true ] && [ false ]`,
			want: [2]bool{false, true},
		},
		{
			val:  `[ true ] && [ -z "" ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ true ] && [ -z "" ] || [ false ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ 2 -gt 1 ] || [ "foo" == "bar" ] && [ false ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ 2 -gt 0 ] || [ "foo" == "foo" ] && [ true ]`,
			want: [2]bool{true, true},
		},
		{
			val:  `[ 2 -lt 0 ] && [ "foo" == "foo" ] || [ true ]`,
			want: [2]bool{true, true},
		},
	}

	for _, c := range cases {
		t.Run(c.val, func(t *testing.T) {
			v := strings.Fields(c.val)
			got, ok := evalConditions(v)
			require.Equal(t, c.want, [2]bool{got, ok})
		})
	}
}
