package variables

import (
	"testing"
)

func TestParseEscapedKeyValue(t *testing.T) {
	var tests = []struct {
		kv string
		k  string
		v  string
		ok bool
	}{
		{"", "", "", false},
		{"=", "", "", true},
		{"key", "key", "", false},
		{"key=", "key", "", true},
		{"key=val", "key", "val", true},
		{"key=val=value=VALUE", "key", "val=value=VALUE", true},
		{"with space=val with space", "with space", "val with space", true},
		{`\==\`, "=", `\`, true},
		{`\===`, "=", `=`, true},
		{`\==\=`, "=", `\=`, true},
		{`value=dmFsdWU=`, "value", `dmFsdWU=`, true},
		{`color\=red=yes!`, "color=red", `yes!`, true},
	}

	for _, tt := range tests {
		k, v, ok := ParseKeyValue(tt.kv)
		Equal(t, tt.k, k)
		Equal(t, tt.v, v)
		Equal(t, tt.ok, ok)
	}
}
