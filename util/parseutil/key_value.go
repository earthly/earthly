package parseutil

import (
	"strings"
)

// ParseArgsForKey parses a list of args and returns the value
// of the first found key
func ParseArgsForKey(key string, args []string) (string, bool) {
	wantNextVal := false
	for _, s := range args {
		if wantNextVal {
			return s, true
		}

		k, v, keyOnly := ParseKeyValue(s)
		if k == key {
			if keyOnly {
				wantNextVal = true
				continue
			}
			return v, true
		}
	}
	return "", false
}

// ParseKeyValue pases a key-value type into its parts
// if a key value needs to contain a = or \, it must be escapped using '\=', and '\\' respectively
// once an unescaped '=' is found, all remaining chars will be used as-is without the need to be escaped.
// the key and value are returned, along with a bool that is true if a value was defined (i.e. an equal was found)
//
// e.g. ParseKeyValue("foo")       -> `foo`,  ``,       false
//      ParseKeyValue("foo=")      -> `foo`,  ``,       true
//      ParseKeyValue("foo=bar")   -> `foo`,  `bar`,    true
//      ParseKeyValue(`f\=oo=bar`) -> `f=oo`, `bar`,    true
//      ParseKeyValue(`foo=bar=`)  -> `foo",  `bar=`,   true
//      ParseKeyValue(`foo=bar\=`) -> `foo",  `bar\=`,  true
func ParseKeyValue(s string) (string, string, bool) {
	key := []string{}
	var escaped bool
	for i, c := range s {
		if escaped {
			key = append(key, string(c))
			escaped = false
			continue
		}
		if c == '\\' {
			escaped = true
			continue
		}
		if c == '=' {
			return strings.Join(key, ""), s[i+1:], true
		}
		key = append(key, string(c))
	}
	return strings.Join(key, ""), "", false
}
