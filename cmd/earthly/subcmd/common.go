package subcmd

import (
	"github.com/urfave/cli/v2"
)

func concatCmds(slices [][]*cli.Command) []*cli.Command {
	var totalLen int

	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]*cli.Command, totalLen)

	var i int

	for _, s := range slices {
		i += copy(result[i:], s)
	}

	return result
}
