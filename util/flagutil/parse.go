package flagutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/jessevdk/go-flags"
	"github.com/urfave/cli/v2"
)

// ArgumentModFunc accepts a flagName which corresponds to the long flag name, and a pointer
// to a flag value. The pointer is nil if no flag was given.
// the function returns a new pointer set to nil if one wants to pretend as if no value was given,
// or a pointer to a new value which will be parsed.
// Note: this was created to allow passing --no-cache=$SOME_VALUE; where we must expand $SOME_VALUE into
// a true/false value before it is parsed. If this feature is used extensively, then it might be time
// to completely fork go-flags with a version where we can include control over expansion struct tags.
type ArgumentModFunc func(flagName string, opt *flags.Option, flagVal *string) (*string, error)

// ParseArgs parses flags and args from a command string
func ParseArgs(command string, data interface{}, args []string) ([]string, error) {
	return ParseArgsWithValueModifier(command, data, args,
		func(_ string, _ *flags.Option, s *string) (*string, error) { return s, nil },
	)
}

// ParseArgsWithValueModifier parses flags and args from a command string; it accepts an optional argumentModFunc
// which is called before each flag value is parsed, and allows one to change the value.
// if the flag value
func ParseArgsWithValueModifier(command string, data interface{}, args []string, argumentModFunc ArgumentModFunc) ([]string, error) {
	return ParseArgsWithValueModifierAndOptions(command, data, args, argumentModFunc, flags.PrintErrors|flags.PassDoubleDash|flags.PassAfterNonOption|flags.AllowBoolValues)
}

// ParseArgsWithValueModifierAndOptions is similar to ParseArgsWithValueModifier, but allows changing the parser options.
func ParseArgsWithValueModifierAndOptions(command string, data interface{}, args []string, argumentModFunc ArgumentModFunc, parserOptions flags.Options) ([]string, error) {
	p := flags.NewNamedParser("", parserOptions)
	var modFuncErr error
	modFunc := func(flagName string, opt *flags.Option, flagVal *string) *string {
		p, err := argumentModFunc(flagName, opt, flagVal)
		if err != nil {
			modFuncErr = err
		}
		return p
	}
	p.ArgumentMod = modFunc
	_, err := p.AddGroup(fmt.Sprintf("%s [options] args", command), "", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initiate parser.AddGroup for %s", command)
	}
	res, err := p.ParseArgs(args)
	if err != nil {
		if parserOptions&flags.PrintErrors != flags.None {
			p.WriteHelp(os.Stderr)
		}
		return nil, err
	}
	if modFuncErr != nil {
		return nil, modFuncErr
	}
	return res, nil
}

// SplitFlagString would return an array of values from the StringSlice, whether it's passed using multiple occuranced of the flag
// or with the values passed with a command.
// For example: --platform linux/amd64 --platform linux/arm64 and --platform "linux/amd64,linux/arm64"
func SplitFlagString(value cli.StringSlice) []string {
	valueStr := strings.TrimLeft(strings.TrimRight(value.String(), "]"), "[")
	return strings.FieldsFunc(valueStr, func(r rune) bool {
		return r == ' ' || r == ','
	})
}
