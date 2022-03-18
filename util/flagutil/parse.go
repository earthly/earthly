package flagutil

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	flags "github.com/jessevdk/go-flags"
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
	p := flags.NewNamedParser("", flags.PrintErrors|flags.PassDoubleDash|flags.PassAfterNonOption|flags.AllowBoolValues)
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
		p.WriteHelp(os.Stderr)
		return nil, err
	}
	if modFuncErr != nil {
		return nil, modFuncErr
	}
	return res, nil
}
