package ast

import "regexp"

var envVarNameRegexp = regexp.MustCompile(`^[a-zA-Z_]+[a-zA-Z0-9_]*$`)

// IsValidEnvVarName returns true if env name is valid
func IsValidEnvVarName(name string) bool {
	return envVarNameRegexp.MatchString(name)
}
