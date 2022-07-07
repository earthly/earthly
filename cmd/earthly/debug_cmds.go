package main

import (
	"encoding/json"
	"fmt"

	"github.com/earthly/earthly/ast"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func (app *earthlyApp) debugCmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "ast",
			Usage:     "Output the AST",
			UsageText: "earthly [options] debug ast",
			Action:    app.actionDebugAst,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "source-map",
					Usage:       "Enable outputting inline sourcemap",
					Destination: &app.enableSourceMap,
				},
			},
		},
	}
}

func (app *earthlyApp) actionDebugAst(cliCtx *cli.Context) error {
	app.commandName = "debugAst"
	if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := "./Earthfile"
	if cliCtx.NArg() == 1 {
		path = cliCtx.Args().First()
	}

	ef, err := ast.Parse(cliCtx.Context, path, app.enableSourceMap)
	if err != nil {
		return err
	}
	efDt, err := json.Marshal(ef)
	if err != nil {
		return errors.Wrap(err, "marshal ast")
	}
	fmt.Print(string(efDt))
	return nil
}
