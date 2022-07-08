package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/containerd/containerd/platforms"
	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/cloud"
	"github.com/moby/buildkit/client"
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
		{
			Name:      "buildkit-info",
			Usage:     "Print the buildkit info",
			UsageText: "earthly [options] debug buildkit-info",
			Action:    app.actionDebugBuildkitInfo,
		},
		{
			Name:      "buildkit-disk-usage",
			Usage:     "Print the buildkit disk usage",
			UsageText: "earthly [options] debug buildkit-disk-usage",
			Action:    app.actionDebugBuildkitDiskUsage,
		},
		{
			Name:      "buildkit-workers",
			Usage:     "Print the buildkit workers",
			UsageText: "earthly [options] debug buildkit-workers",
			Action:    app.actionDebugBuildkitWorkers,
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

func (app *earthlyApp) actionDebugBuildkitInfo(cliCtx *cli.Context) error {
	app.commandName = "debugBuildkitInfo"

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	bkClient, err := app.getBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()

	info, err := bkClient.Info(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "get buildkit info")
	}

	fmt.Printf("Buildkit version: %s\n", info.BuildkitVersion.Version)
	fmt.Printf("Buildkit revision: %s\n", info.BuildkitVersion.Revision)
	fmt.Printf("Buildkit package: %s\n", info.BuildkitVersion.Package)
	fmt.Printf("Num sessions: %d\n", info.NumSessions)
	return nil
}

func (app *earthlyApp) actionDebugBuildkitDiskUsage(cliCtx *cli.Context) error {
	app.commandName = "debugBuildkitDiskUsage"

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	bkClient, err := app.getBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()

	infos, err := bkClient.DiskUsage(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "get buildkit disk usage")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "ID\tDescription\tSize\tUsageCount\tRecordType\tMutable\tShared\tInUse\tParents\tCreatedAt\tLastUsedAt\n")
	for _, info := range infos {
		var rt string
		switch info.RecordType {
		case client.UsageRecordTypeCacheMount:
			rt = "cache"
		case client.UsageRecordTypeFrontend:
			rt = "frontend"
		case client.UsageRecordTypeGitCheckout:
			rt = "git"
		case client.UsageRecordTypeInternal:
			rt = "internal"
		case client.UsageRecordTypeLocalSource:
			rt = "local-source"
		case client.UsageRecordTypeRegular:
			rt = "regular"
		}
		fmt.Fprintf(
			w, "%s\t%s\t%d\t%d\t%s\t%t\t%t\t%t\t%s\t%s\t%s\n",
			info.ID, info.Description, info.Size, info.UsageCount,
			rt, info.Mutable, info.Shared, info.InUse,
			strings.Join(info.Parents, ","),
			info.CreatedAt.Format(time.RFC3339), info.LastUsedAt.Format(time.RFC3339),
		)
	}
	return nil
}

func (app *earthlyApp) actionDebugBuildkitWorkers(cliCtx *cli.Context) error {
	app.commandName = "debugBuildkitWorkers"

	cloudClient, err := cloud.NewClient(app.apiServer, app.sshAuthSock, app.authToken, app.console.Warnf)
	if err != nil {
		return errors.Wrap(err, "failed to create cloud client")
	}
	bkClient, err := app.getBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()

	workers, err := bkClient.ListWorkers(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "get buildkit workers")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()
	fmt.Fprintf(w, "ID\tLabels\tPlatforms\tVersion\tPackage\tRevision\tParallelism current\tParallelism max\tParallelism waiting\n")
	for _, info := range workers {
		ps := make([]string, 0, len(info.Platforms))
		for _, p := range info.Platforms {
			ps = append(ps, platforms.Format(p))
		}
		ls := make([]string, 0, len(info.Labels))
		for lk, lv := range info.Labels {
			ls = append(ls, fmt.Sprintf("%s=%s", lk, lv))
		}

		fmt.Fprintf(
			w, "%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\n",
			info.ID, strings.Join(ls, ","), strings.Join(ps, ","),
			info.BuildkitVersion.Version, info.BuildkitVersion.Package,
			info.BuildkitVersion.Revision,
			info.ParallelismCurrent, info.ParallelismMax,
			info.ParallelismWaiting)
	}
	return nil
}
