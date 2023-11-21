package inputgraph

import (
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
