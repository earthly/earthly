package variables

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// ParseKeyValue pases a key-value type into its parts
// if a key value needs to contain a = or \, it must be escaped using '\=', and '\\' respectively
// once an unescaped '=' is found, all remaining chars will be used as-is without the need to be escaped.
// the key and value are returned, along with a bool that is true if a value was defined (i.e. an equal was found)
//
// e.g. ParseKeyValue("foo")       -> `foo`,  “,       false
//
//	ParseKeyValue("foo=")      -> `foo`,  ``,       true
//	ParseKeyValue("foo=bar")   -> `foo`,  `bar`,    true
//	ParseKeyValue(`f\=oo=bar`) -> `f=oo`, `bar`,    true
//	ParseKeyValue(`foo=bar=`)  -> `foo",  `bar=`,   true
//	ParseKeyValue(`foo=bar\=`) -> `foo",  `bar\=`,  true
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

// AddEnv takes in a slice of env vars in key-value format and a new key-value
// string to it, taking care of possible overrides.
func AddEnv(envVars []string, key, value string) []string {
	// Note that this mutates the original slice.
	envMap, err := godotenv.Unmarshal(fmt.Sprintf("%s=%s", key, value))

	if err != nil {
		logrus.Errorf("Error adding env: due to %v", err)
		return envVars
	}

	for k, v := range envMap {
		env := fmt.Sprintf("%s=%s", k, v)
		logrus.Debugf("Adding env: %s", env)
		envVars = append(envVars, env)
	}
	return envVars
}
