package shell

import "github.com/earthly/earthly/ast"

// IsValidEnvVarName returns true if env name is valid
func IsValidEnvVarName(name string) bool {
	return ast.IsValidEnvVarName(name)
}
