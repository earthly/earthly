package variables

import (
	"fmt"
	"strings"

	"github.com/earthly/earthly/util/shell"
	"github.com/pkg/errors"
)

// ParseFlagArgs parses flag-form args.
// These can be represented as `--arg=value` or `--arg value`.
// The result is a slice that can be passed into ParseArgs or to ParseCommandLineArgs.
func ParseFlagArgs(args []string) ([]string, error) {
	flags, nonFlags, err := ParseFlagArgsWithNonFlags(args)
	if err != nil {
		return nil, err
	}
	if len(nonFlags) != 0 {
		return nil, errors.Errorf("invalid argument %s", nonFlags[0])
	}
	return flags, nil
}

// ParseFlagArgsWithNonFlags parses flag-form args together with any possible optional additional
// args. e.g. "--flag1=value arg1 --flag2=value --flag3=value arg2 arg3".
func ParseFlagArgsWithNonFlags(args []string) ([]string, []string, error) {
	flags := make([]string, 0, len(args))
	nonFlags := []string{}
	keyFromPrev := ""
	for _, arg := range args {
		var k, v string
		if keyFromPrev != "" {
			k = keyFromPrev
			keyFromPrev = ""
			v = arg
		} else {
			var trimmedArg string
			if strings.HasPrefix(arg, "--") {
				trimmedArg = strings.TrimPrefix(arg, "--")
			} else if strings.HasPrefix(arg, "-") {
				trimmedArg = strings.TrimPrefix(arg, "-")
			} else {
				nonFlags = append(nonFlags, arg)
				continue
			}
			var hasValue bool
			k, v, hasValue = ParseKeyValue(trimmedArg)
			if !shell.IsValidEnvVarName(k) {
				return nil, nil, errors.Errorf("invalid arg name: %s", arg)
			}
			if !hasValue {
				keyFromPrev = k
				continue
			}
		}
		escK := strings.ReplaceAll(k, "=", "\\=")
		flags = append(flags, fmt.Sprintf("%s=%s", escK, v))
	}
	if keyFromPrev != "" {
		return nil, nil, errors.Errorf("no value provided for --%s", keyFromPrev)
	}
	return flags, nonFlags, nil
}
