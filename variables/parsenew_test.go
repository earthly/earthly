package variables

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestParseFlagArgs(t *testing.T) {
	var tests = []struct {
		kvFlag []string
		kv     []string
	}{
		{[]string{}, []string{}},
		{[]string{"--flag=foo"}, []string{"flag=foo"}},
		{[]string{"--flag", "foo"}, []string{"flag=foo"}},
		{[]string{"--flag", "--foo"}, []string{"flag=--foo"}},
		{[]string{"--flag=--foo"}, []string{"flag=--foo"}},
		{[]string{"--flag\\=name=foo"}, []string{"flag\\=name=foo"}},
		{[]string{"--flag\\=name=foo\\=bar"}, []string{"flag\\=name=foo\\=bar"}},
		{[]string{"--flag\\=name=foo=bar"}, []string{"flag\\=name=foo=bar"}},
		{[]string{"--flag=foo", "--another=bar"}, []string{"flag=foo", "another=bar"}},
		{[]string{"--flag", "--foo", "--another", "--bar"}, []string{"flag=--foo", "another=--bar"}},
	}

	for _, tt := range tests {
		kvs, err := ParseFlagArgs(tt.kvFlag)
		Equal(t, kvs, tt.kv)
		NoError(t, err)
	}
}

func TestNegativeParseFlagArgs(t *testing.T) {
	var tests = []struct {
		kvFlag []string
	}{
		{[]string{"--foo"}},
		{[]string{"--foo", "--bar", "--bar2"}},
		{[]string{"-foo", "bar"}},
		{[]string{"-foo=bar"}},
		{[]string{"foo=bar"}},
	}

	for _, tt := range tests {
		_, err := ParseFlagArgs(tt.kvFlag)
		Error(t, err)
	}
}
