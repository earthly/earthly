package subcmd

import (
	"context"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/earthly/earthly/ast/hint"

	"github.com/earthly/earthly/util/proj"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

const efIndent = "    "

type Init struct {
	cli CLI
}

func NewInit(cli CLI) *Init {
	return &Init{
		cli: cli,
	}
}

func (a *Init) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "init",
			Description: "*experimental* Initialize a project.",
			Usage:       "*experimental* Initialize an Earthfile for the current project",
			Action:      a.action,
		},
	}
}

func (a *Init) action(cliCtx *cli.Context) error {
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
	if len(projs) > 1 {
		// This is easy enough to support when we have more than one project
		// type, but for now there's no point.
		return errors.Errorf("%d projects detected, but multiple project types are not supported by init yet", len(projs))
	}

	p := projs[0]
	if p.Root(ctx) != absWd {
		// In the distant future, this may be used to generate multiple
		// Earthfiles over multiple directories and call them from a main
		// Earthfile target with BUILD.
		return errors.Errorf("project type %T wants to generate an Earthfile in an unsupported directory: %q", p, p.Root(ctx))
	}

	return initSingleProject(ctx, f, p)
}

func initSingleProject(ctx context.Context, w io.Writer, p proj.Project) error {
	tgts, err := p.Targets(ctx)
	if err != nil {
		return errors.Wrapf(err, "could not generate targets for project type %T", p)
	}
	for i, tgt := range tgts {
		tgt.SetPrefix(ctx, "")
		if i > 0 {
			_, err = w.Write([]byte("\n"))
			if err != nil {
				return errors.Wrapf(err, "could not write newline separator between targets")
			}
		}
		err := tgt.Format(ctx, w, efIndent, 0)
		if err != nil {
			return errors.Wrapf(err, "could not format target for project type %T", p)
		}
	}
	return nil
}
