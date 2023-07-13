package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/earthly/earthly/ast/hint"
	"github.com/earthly/earthly/util/proj"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const efIndent = "    "

func (app *earthlyApp) actionInit(cliCtx *cli.Context) error {
	ctx := cliCtx.Context

	wd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "could not load current working directory")
	}
	absWd, err := filepath.Abs(wd)
	if err != nil {
		return errors.Wrapf(err, "could not get absolute path for %q", wd)
	}

	efPath := filepath.Join(absWd, "Earthfile")
	_, err = os.Stat(efPath)
	if err == nil {
		return hint.Wrap(fs.ErrExist, "an Earthfile already exists; if you want to re-init the project, remove the Earthfile first.")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		return errors.Wrap(err, "could not check for existing Earthfile")
	}

	projs, err := proj.All(ctx, absWd)
	if err != nil {
		return errors.Wrapf(err, "could not get projects for %q", absWd)
	}
	if len(projs) == 0 {
		return errors.Errorf("no supported projects found in directory %q", absWd)
	}

	f, err := os.Create(efPath)
	if err != nil {
		return errors.Wrapf(err, "could not create %q", efPath)
	}
	defer f.Close()

	_, err = f.WriteString("VERSION --arg-scope-and-set 0.7\n\n")
	if err != nil {
		return errors.Wrapf(err, "could not write version string in %q", efPath)
	}
	for _, p := range projs {
		if p.Root(ctx) != absWd {
			return errors.Errorf("project type %T wants to generate an Earthfile in an unsupported directory: %q", p, p.Root(ctx))
		}
		tgts, err := p.Targets(ctx)
		if err != nil {
			return errors.Wrapf(err, "could not generate targets for project type %T", p)
		}
		for _, tgt := range tgts {
			err := tgt.Format(ctx, f, efIndent, 0)
			if err != nil {
				return errors.Wrapf(err, "could not format target for project type %T", p)
			}
			_, err = f.WriteString("\n")
			if err != nil {
				return errors.Wrapf(err, "could not write newline separator between targets")
			}
		}
	}
	return nil
}
