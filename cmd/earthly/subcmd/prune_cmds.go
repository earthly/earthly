package subcmd

import (
	"time"

	"github.com/dustin/go-humanize"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/earthly/earthly/util/flagutil"
)

type Prune struct {
	cli CLI

	all          bool
	reset        bool
	keepDuration flagutil.Duration
	targetSize   flagutil.ByteSizeValue
}

func NewPrune(cli CLI) *Prune {
	return &Prune{
		cli: cli,
	}
}

func (a *Prune) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:  "prune",
			Usage: "Prune Earthly build cache",
			Description: `Prune Earthly build cache in one of two forms.
	Standard Form:
		Issues a prune command on the BuildKit daemon.
	Reset Form:
		Restarts the BuildKit daemon and instructs it to complete delete the cache
		directory on startup.`,
			Action: a.action,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "all",
					Aliases:     []string{"a"},
					EnvVars:     []string{"EARTHLY_PRUNE_ALL"},
					Usage:       "Prune all cache via BuildKit daemon",
					Destination: &a.all,
				},
				&cli.BoolFlag{
					Name:    "reset",
					EnvVars: []string{"EARTHLY_PRUNE_RESET"},
					Usage: `Reset cache entirely by restarting BuildKit daemon and wiping cache dir.
				This option is not available when using satellites.`,
					Destination: &a.reset,
				},
				&cli.GenericFlag{
					Name: "age",
					Usage: `Prune cache older than the specified duration passed in as a string;
						duration is specified with an integer value followed by a m, h, or d suffix which represents minutes, hours, or days respectively, e.g. 24h, or 1d`,
					Value: &a.keepDuration,
				},
				&cli.GenericFlag{
					Name:  "size",
					Usage: "Prune cache to specified size, starting from oldest",
					Value: &a.targetSize,
				},
			},
		},
	}
}

func (a *Prune) action(cliCtx *cli.Context) error {
	a.cli.SetCommandName("prune")
	if cliCtx.NArg() != 0 {
		return errors.New("invalid arguments")
	}
	if a.reset {
		if a.cli.IsUsingSatellite(cliCtx) {
			return errors.New("Cannot prune --reset when using a satellite. Try without --reset")
		}
		err := a.cli.InitFrontend(cliCtx)
		if err != nil {
			return err
		}
		err = buildkitd.ResetCache(cliCtx.Context, a.cli.Console(), a.cli.Flags().BuildkitdImage, a.cli.Flags().ContainerName, a.cli.Flags().InstallationName, a.cli.Flags().ContainerFrontend, a.cli.Flags().BuildkitdSettings)
		if err != nil {
			return errors.Wrap(err, "reset cache")
		}
		return nil
	}

	// Prune via API.
	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	bkClient, cleanupTLS, err := a.cli.GetBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "prune new buildkitd client")
	}
	defer bkClient.Close()
	defer cleanupTLS()
	var opts []client.PruneOption

	if a.all {
		opts = append(opts, client.PruneAll)
	}

	if a.keepDuration > 0 || a.targetSize > 0 {
		opts = append(opts, client.WithKeepOpt(time.Duration(a.keepDuration), int64(a.targetSize)))
	}

	ch := make(chan client.UsageInfo, 1)
	eg, ctx := errgroup.WithContext(cliCtx.Context)
	eg.Go(func() error {
		err = bkClient.Prune(ctx, ch, opts...)
		if err != nil {
			return errors.Wrap(err, "buildkit prune")
		}
		close(ch)
		return nil
	})

	total := uint64(0)
	eg.Go(func() error {
		for {
			select {
			case usageInfo, ok := <-ch:
				if !ok {
					return nil
				}
				a.cli.Console().Printf("%s\t%s\n", usageInfo.ID, humanize.Bytes(uint64(usageInfo.Size)))
				total += uint64(usageInfo.Size)
			case <-ctx.Done():
				return nil
			}
		}
	})
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "err group")
	}
	a.cli.Console().Printf("Freed %s\n", humanize.Bytes(total))
	return nil
}
