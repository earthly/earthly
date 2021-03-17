package variables

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ParseFlagArgs parses flag-form args.
// These can be represented as `--arg=value` or `--arg value`.
// The result is a slice that can be passed in to ParseArgs or to ParseCommandLineArgs.
func ParseFlagArgs(args []string) ([]string, error) {
	ret := make([]string, 0, len(args))
	keyFromPrev := ""
	for _, arg := range args {
		var k, v string
		if keyFromPrev != "" {
			k = keyFromPrev
			keyFromPrev = ""
			v = arg
		} else {
			if !strings.HasPrefix(arg, "--") {
				return nil, errors.Errorf("invalid flag %s", arg)
			}
			trimmedArg := strings.TrimPrefix(arg, "--")
			var hasValue bool
			k, v, hasValue = ParseKeyValue(trimmedArg)
			if !hasValue {
				keyFromPrev = k
				continue
			}
		}
		escK := strings.ReplaceAll(k, "=", "\\=")
		ret = append(ret, fmt.Sprintf("%s=%s", escK, v))
	}
	if keyFromPrev != "" {
		return nil, errors.Errorf("no value provided for --%s", keyFromPrev)
	}
	return ret, nil
}
