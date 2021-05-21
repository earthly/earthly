package flagutil

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	flags "github.com/jessevdk/go-flags"
)

// ParseArgs parses flags and args from a command string
func ParseArgs(command string, data interface{}, args []string) ([]string, error) {
	p := flags.NewNamedParser("", flags.PrintErrors|flags.PassDoubleDash|flags.PassAfterNonOption)
	_, err := p.AddGroup(fmt.Sprintf("%s [options] args", command), "", data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initiate parser.AddGroup for %s", command)
	}
	res, err := p.ParseArgs(args)
	if err != nil {
		p.WriteHelp(os.Stderr)
		return nil, err
	}
	return res, nil
}
