package subcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/containerd/containerd/platforms"
	"github.com/dustin/go-humanize"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/earthly/earthly/ast"
	"github.com/earthly/earthly/cmd/earthly/helper"
)

type Debug struct {
	cli CLI

	enableSourceMap bool
}

func NewDebug(cli CLI) *Debug {
	return &Debug{
		cli: cli,
	}
}

func (a *Debug) Cmds() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "debug",
			Usage:       "Print debug information about an Earthfile",
			Description: "Print debug information about an Earthfile.",
			ArgsUsage:   "[<path>]",
			Hidden:      true, // Dev purposes only.
			Subcommands: []*cli.Command{
				{
					Name:        "ast",
					Usage:       "Output the AST",
					UsageText:   "earthly [options] debug ast",
					Description: "Output the AST.",
					Action:      a.actionAst,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:        "source-map",
							Usage:       "Enable outputting inline sourcemap",
							Destination: &a.enableSourceMap,
						},
					},
				},
				{
					Name:        "buildkit-info",
					Usage:       "Print the buildkit info",
					UsageText:   "earthly [options] debug buildkit-info",
					Description: "Print the builtkit info.",
					Action:      a.actionBuildkitInfo,
				},
				{
					Name:        "buildkit-disk-usage",
					Usage:       "Print the buildkit disk usage",
					UsageText:   "earthly [options] debug buildkit-disk-usage",
					Description: "Print the buildkit disk usage.",
					Action:      a.actionBuildkitDiskUsage,
				},
				{
					Name:        "buildkit-workers",
					Usage:       "Print the buildkit workers",
					UsageText:   "earthly [options] debug buildkit-workers",
					Description: "Print the buildkit workers.",
					Action:      a.actionBuildkitWorkers,
				},
				{
					Name:        "buildkit-shutdown-if-idle",
					Usage:       "Shutdown the buildkit if it is idle",
					UsageText:   "earthly [options] debug buildkit-shutdown-if-idle",
					Description: "Shutdown the buildkit if it is idle.",
					Action:      a.actionBuildkitShutdownIfIdle,
				},
				{
					Name:        "buildkit-session-history",
					Usage:       "Print the buildkit session history",
					UsageText:   "earthly [options] debug buildkit-session-history",
					Description: "Print the buildkit session history.",
					Action:      a.actionBuildkitSessionHistory,
				},
			},
		},
	}
}

func (a *Debug) actionAst(cliCtx *cli.Context) error {
	a.cli.SetCommandName("debugAst")

	if cliCtx.NArg() > 1 {
		return errors.New("invalid number of arguments provided")
	}
	path := "./Earthfile"
	if cliCtx.NArg() == 1 {
		path = cliCtx.Args().First()
	}

	ef, err := ast.Parse(cliCtx.Context, path, a.enableSourceMap)
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

func (a *Debug) actionBuildkitSessionHistory(cliCtx *cli.Context) error {
	a.cli.SetCommandName("debugBuildkitSessions")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	bkClient, cleanupTLS, err := a.cli.GetBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	defer cleanupTLS()

	history, err := bkClient.SessionHistory(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "get buildkit session history")
	}

	byt, _ := json.MarshalIndent(history, "", "  ")
	fmt.Println(string(byt))
	return nil
}

func (a *Debug) actionBuildkitInfo(cliCtx *cli.Context) error {
	a.cli.SetCommandName("debugBuildkitInfo")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	bkClient, cleanupTLS, err := a.cli.GetBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	defer cleanupTLS()

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

func (a *Debug) actionBuildkitDiskUsage(cliCtx *cli.Context) error {
	a.cli.SetCommandName("debugBuildkitDiskUsage")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	bkClient, cleanupTLS, err := a.cli.GetBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	defer cleanupTLS()

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
		lua := "nil"
		if info.LastUsedAt != nil {
			lua = info.LastUsedAt.Format(time.RFC3339)
		}
		fmt.Fprintf(
			w, "%s\t%s\t%d\t%d\t%s\t%t\t%t\t%t\t%s\t%s\t%s\n",
			info.ID, info.Description, info.Size, info.UsageCount,
			rt, info.Mutable, info.Shared, info.InUse,
			strings.Join(info.Parents, ","),
			info.CreatedAt.Format(time.RFC3339), lua,
		)
	}
	return nil
}

func (a *Debug) actionBuildkitWorkers(cliCtx *cli.Context) error {
	a.cli.SetCommandName("debugBuildkitWorkers")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	bkClient, cleanupTLS, err := a.cli.GetBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	defer cleanupTLS()

	workers, err := bkClient.ListWorkers(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "get buildkit workers")
	}

	for _, info := range workers {
		ps := make([]string, 0, len(info.Platforms))
		for _, p := range info.Platforms {
			ps = append(ps, platforms.Format(p))
		}
		ls := make([]string, 0, len(info.Labels))
		for lk, lv := range info.Labels {
			ls = append(ls, fmt.Sprintf("%s=%s", lk, lv))
		}
		fmt.Printf("Worker %s\n", info.ID)
		fmt.Printf("\tLabels: %s\n", strings.Join(ls, ","))
		fmt.Printf("\tPlatforms: %s\n", strings.Join(ps, ","))
		fmt.Printf("\tVersion: %s\n", info.BuildkitVersion.Version)
		fmt.Printf("\tPackage: %s\n", info.BuildkitVersion.Package)
		fmt.Printf("\tRevision: %s\n", info.BuildkitVersion.Revision)

		fmt.Printf("\tParallelism current: %d\n", info.ParallelismCurrent)
		fmt.Printf("\tParallelism max: %d\n", info.ParallelismMax)
		fmt.Printf("\tParallelism waiting: %d\n", info.ParallelismWaiting)

		fmt.Printf("\tGC Num runs summary: %d\n", info.GCAnalytics.NumRuns)
		fmt.Printf("\tGC Num failures: %d\n", info.GCAnalytics.NumFailures)
		fmt.Printf("\tGC Avg duration: %s\n", info.GCAnalytics.AvgDuration)
		fmt.Printf("\tGC Avg records cleared: %d\n", info.GCAnalytics.AvgRecordsCleared)
		fmt.Printf("\tGC Avg size cleared: %s\n", humanize.Bytes(uint64(info.GCAnalytics.AvgSizeCleared)))
		fmt.Printf("\tGC Avg records before: %d\n", info.GCAnalytics.AvgRecordsBefore)
		fmt.Printf("\tGC Avg size before: %s\n", humanize.Bytes(uint64(info.GCAnalytics.AvgSizeBefore)))
		fmt.Printf("\tGC All-time runs: %d\n", info.GCAnalytics.AllTimeRuns)
		fmt.Printf("\tGC All-time max duration: %s\n", humanizeDuration(info.GCAnalytics.AllTimeMaxDuration))
		fmt.Printf("\tGC All-time duration: %s\n", humanizeDuration(info.GCAnalytics.AllTimeDuration))
		if info.GCAnalytics.CurrentStartTime != nil {
			fmt.Printf("\tGC Current start time: %s\n", humanizeTime(info.GCAnalytics.CurrentStartTime))
			fmt.Printf("\tGC Current num records before: %d\n", info.GCAnalytics.CurrentNumRecordsBefore)
			fmt.Printf("\tGC Current size before: %s\n", humanize.Bytes(uint64(info.GCAnalytics.CurrentSizeBefore)))
		} else {
			fmt.Printf("\tNo GC run currently ongoing\n")
		}
		if info.GCAnalytics.LastStartTime != nil {
			fmt.Printf("\tGC Last start time: %s\n", humanizeTime(info.GCAnalytics.LastStartTime))
			fmt.Printf("\tGC Last end time: %s\n", humanizeTime(info.GCAnalytics.LastEndTime))
			fmt.Printf(
				"\tGC Last duration: %s\n",
				humanizeDuration(info.GCAnalytics.LastEndTime.Sub(*info.GCAnalytics.LastStartTime)))
			fmt.Printf("\tGC Last num records before: %d\n", info.GCAnalytics.LastNumRecordsBefore)
			fmt.Printf("\tGC Last size before: %s\n", humanize.Bytes(uint64(info.GCAnalytics.LastSizeBefore)))
			fmt.Printf("\tGC Last num records cleared: %d\n", info.GCAnalytics.LastNumRecordsCleared)
			fmt.Printf("\tGC Last size cleared: %s\n", humanize.Bytes(uint64(info.GCAnalytics.LastSizeCleared)))
			fmt.Printf("\tGC Last success: %v\n", info.GCAnalytics.LastSuccess)
		} else {
			fmt.Printf("\tGC has not run yet\n")
		}
	}
	return nil
}

func (a *Debug) actionBuildkitShutdownIfIdle(cliCtx *cli.Context) error {
	a.cli.SetCommandName("debugBuildkitShutdownIfIdle")

	cloudClient, err := helper.NewCloudClient(a.cli)
	if err != nil {
		return err
	}
	bkClient, cleanupTLS, err := a.cli.GetBuildkitClient(cliCtx, cloudClient)
	if err != nil {
		return errors.Wrap(err, "build new buildkitd client")
	}
	defer bkClient.Close()
	defer cleanupTLS()

	ok, numSessions, err := bkClient.ShutdownIfIdle(cliCtx.Context)
	if err != nil {
		return errors.Wrap(err, "shutdown buildkit if idle")
	}
	fmt.Printf("Shutting down: %t\n", ok)
	fmt.Printf("Num sessions: %d\n", numSessions)
	return nil
}

func humanizeDuration(d time.Duration) string {
	return fmt.Sprintf("%v", d.Round(time.Second))
}

func humanizeTime(t *time.Time) string {
	if t == nil {
		return "nil"
	}
	if t.IsZero() {
		return "zero"
	}
	return t.Format(time.RFC3339)
}
