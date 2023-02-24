package variables

import (
	"testing"
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
		{[]string{"-flag", "--foo", "-another", "bar"}, []string{"flag=--foo", "another=bar"}},
	}

	for _, tt := range tests {
		kvs, err := ParseFlagArgs(tt.kvFlag)
		NoError(t, err)
		Equal(t, kvs, tt.kv)
	}
}

func TestNegativeParseFlagArgs(t *testing.T) {
	var tests = []struct {
		kvFlag []string
	}{
		{[]string{"--foo"}},
		{[]string{"--foo", "--bar", "--bar2"}},
		{[]string{"foo=bar"}},
	}

	for _, tt := range tests {
		_, err := ParseFlagArgs(tt.kvFlag)
		Error(t, err)
	}
}

func TestParseFlagArgsWithNonFlags(t *testing.T) {
	var tests = []struct {
		kvFlag   []string
		flags    []string
		nonFlags []string
	}{
		{[]string{}, []string{}, []string{}},
		{[]string{"--flag=foo"}, []string{"flag=foo"}, []string{}},
		{[]string{"--flag=foo", "arg"}, []string{"flag=foo"}, []string{"arg"}},
		{[]string{"arg", "--flag=foo"}, []string{"flag=foo"}, []string{"arg"}},
		{[]string{"arg", "--flag=foo", "arg2", "--flag2=bar"}, []string{"flag=foo", "flag2=bar"}, []string{"arg", "arg2"}},
		{[]string{"arg", "--flag", "foo", "arg2", "--flag2=bar"}, []string{"flag=foo", "flag2=bar"}, []string{"arg", "arg2"}},
		{[]string{"arg", "-flag", "foo", "arg2", "-flag2=bar"}, []string{"flag=foo", "flag2=bar"}, []string{"arg", "arg2"}},
		{[]string{"arg"}, []string{}, []string{"arg"}},
		{[]string{"just", "args"}, []string{}, []string{"just", "args"}},
	}

	for _, tt := range tests {
		flags, nonFlags, err := ParseFlagArgsWithNonFlags(tt.kvFlag)
		NoError(t, err)
		Equal(t, flags, tt.flags)
		Equal(t, nonFlags, tt.nonFlags)
	}
}
