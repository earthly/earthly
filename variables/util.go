package variables

import (
	"fmt"

	"github.com/earthly/earthly/util/parseutil"
)

// AddEnv takes in a slice of env vars in key-value format and a new key-value
// string to it, taking care of possible overrides.
func AddEnv(envVars []string, key, value string) []string {
	// Note that this mutates the original slice.
	found := false
	for i, envVar := range envVars {
		k, _, _ := parseutil.ParseKeyValue(envVar)
		if k == key {
			envVars[i] = fmt.Sprintf("%s=%s", key, value)
			found = true
			break
		}
	}
	if !found {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}
	return envVars
}
