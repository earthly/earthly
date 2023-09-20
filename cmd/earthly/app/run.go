package app

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/earthly/earthly/cloud"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/cmd/earthly/helper"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/analytics"
	"github.com/earthly/earthly/builder"
	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/errutil"
	"github.com/earthly/earthly/util/reflectutil"
)

var runExitCodeRegexp = regexp.MustCompile(`did not complete successfully: exit code: [^0][0-9]*$`)

func (app *EarthlyApp) Run(ctx context.Context, console conslogging.ConsoleLogger, startTime time.Time, lastSignal os.Signal) int {
	err := app.unhideFlags(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error un-hiding flags %v", err)
		os.Exit(1)
	}
	helper.AutoComplete(ctx, app.BaseCLI)

	exitCode := app.run(ctx, os.Args)

	// app.Cfg will be nil when a user runs `earthly --version`;
	// however in all other regular commands app.Cfg will be set in app.Before
	if !app.BaseCLI.Flags().DisableAnalytics && app.BaseCLI.Cfg() != nil && !app.BaseCLI.Cfg().Global.DisableAnalytics {
		// Use a new context, in case the original context is cancelled due to sigint.
		ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		displayErrors := app.BaseCLI.Flags().Verbose
		cloudClient, err := helper.NewCloudClient(app.BaseCLI)
		if err != nil && displayErrors {
			app.BaseCLI.Console().Warnf("unable to start cloud app.BaseCLIent: %s", err)
		} else if err == nil {
			analytics.AddCLIProject(app.BaseCLI.Flags().OrgName, app.BaseCLI.Flags().ProjectName)
			org, project := analytics.ProjectDetails()
			analytics.CollectAnalytics(
				ctxTimeout, cloudClient, displayErrors, analytics.Meta{
					Version:          app.BaseCLI.Version(),
					Platform:         common.GetPlatform(),
					BuildkitPlatform: app.BaseCLI.AnaMetaBKPlatform(),
					UserPlatform:     app.BaseCLI.AnaMetaUserPlatform(),
					GitSHA:           app.BaseCLI.GitSHA(),
					CommandName:      app.BaseCLI.CommandName(),
					ExitCode:         exitCode,
					Target:           app.BaseCLI.AnaMetaTarget(),
					IsSatellite:      app.BaseCLI.AnaMetaIsSat(),
					SatelliteVersion: app.BaseCLI.AnaMetaSatCurrentVersion(),
					IsRemoteBuildkit: app.BaseCLI.AnaMetaIsRemoteBK(),
					Realtime:         time.Since(startTime),
					OrgName:          org,
					ProjectName:      project,
					EarthlyCIRunner:  app.BaseCLI.Flags().EarthlyCIRunner,
				},
				app.BaseCLI.Flags().InstallationName,
			)
		}
	}
	if lastSignal != nil {
		app.BaseCLI.Console().PrintBar(color.New(color.FgHiYellow), fmt.Sprintf("WARNING: exiting due to received %s signal", lastSignal.String()), "")
	}

	return exitCode
}

func (app *EarthlyApp) unhideFlags(ctx context.Context) error {
	var err error
	if os.Getenv("EARTHLY_AUTOCOMPLETE_HIDDEN") != "" && os.Getenv("COMP_POINT") == "" { // TODO delete this check after 2022-03-01
		// only display warning when NOT under complete mode (otherwise we break auto completion)
		app.BaseCLI.Console().Warnf("Warning: EARTHLY_AUTOCOMPLETE_HIDDEN has been renamed to EARTHLY_SHOW_HIDDEN\n")
	}
	showHidden := false
	showHiddenStr := os.Getenv("EARTHLY_SHOW_HIDDEN")
	if showHiddenStr != "" {
		showHidden, err = strconv.ParseBool(showHiddenStr)
		if err != nil {
			return err
		}
	}
	if !showHidden {
		return nil
	}

	for _, fl := range app.BaseCLI.App().Flags {
		reflectutil.SetBool(fl, "Hidden", false)
	}

	unhideFlagsCommands(ctx, app.BaseCLI.App().Commands)

	return nil
}

func unhideFlagsCommands(ctx context.Context, cmds []*cli.Command) {
	for _, cmd := range cmds {
		reflectutil.SetBool(cmd, "Hidden", false)
		for _, flg := range cmd.Flags {
			reflectutil.SetBool(flg, "Hidden", false)
		}
		unhideFlagsCommands(ctx, cmd.Subcommands)
	}
}

func (app *EarthlyApp) run(ctx context.Context, args []string) int {
	defer func() {
		if app.BaseCLI.LogbusSetup() != nil {
			err := app.BaseCLI.LogbusSetup().Close()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error(s) in logbus: %v", err)
			}
			if app.BaseCLI.Flags().LogstreamDebugManifestFile != "" {
				err := app.BaseCLI.LogbusSetup().DumpManifestToFile(app.BaseCLI.Flags().LogstreamDebugManifestFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error dumping manifest: %v", err)
				}
			}
		}
	}()
	app.BaseCLI.Logbus().Run().SetStart(time.Now())
	// Initialize log streaming early if we're passed the organization and
	// project names as environmental variables. This will allow nearly all
	// initialization errors to be surfaced to the log streaming service. Access
	// to this organization and project will be verified when the stream begins.

	if app.BaseCLI.Flags().OrgName != "" && app.BaseCLI.Flags().ProjectName != "" && app.BaseCLI.Cfg() != nil && !app.BaseCLI.Cfg().Global.DisableLogSharing && app.BaseCLI.Flags().LogstreamUpload {
		cloudClient, err := helper.NewCloudClient(app.BaseCLI, cloud.WithLogstreamGRPCAddressOverride(app.BaseCLI.Flags().LogstreamAddressOverride))
		if err != nil {
			app.BaseCLI.Console().Warnf("Failed to initialize cloud app.BaseCLIent: %v", err)
			return 1
		}
		if cloudClient.IsLoggedIn(ctx) {
			app.BaseCLI.Console().VerbosePrintf("Logbus: setting organization %q and project %q", app.BaseCLI.Flags().OrgName, app.BaseCLI.Flags().ProjectName)
			analytics.AddEarthfileProject(app.BaseCLI.Flags().OrgName, app.BaseCLI.Flags().ProjectName)
			app.BaseCLI.LogbusSetup().SetOrgAndProject(app.BaseCLI.Flags().OrgName, app.BaseCLI.Flags().ProjectName)
			app.BaseCLI.LogbusSetup().StartLogStreamer(ctx, cloudClient)
			logstreamURL := fmt.Sprintf("%s/builds/%s", app.BaseCLI.CIHost(), app.BaseCLI.LogbusSetup().InitialManifest.GetBuildId())
			app.BaseCLI.Console().ColorPrintf(color.New(color.FgHiYellow), "Streaming logs to %s\n\n", logstreamURL)
		}
	}

	defer func() {
		// Just in case this is forgotten somewhere else.
		app.BaseCLI.Logbus().Run().SetFatalError(
			time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_OTHER,
			"No SetFatalError called appropriately. This should never happen.")
	}()
	rpcRegex := regexp.MustCompile(`(?U)rpc error: code = .+ desc = `)

	err := app.BaseCLI.App().RunContext(ctx, args)
	if err != nil {
		ie, isInterpreterError := earthfile2llb.GetInterpreterError(err)
		var failedOutput string
		var buildErr *builder.BuildError
		if errors.As(err, &buildErr) {
			failedOutput = buildErr.VertexLog()
		}
		if app.BaseCLI.Flags().Debug {
			// Get the stack trace from the deepest error that has it and print it.
			type stackTracer interface {
				StackTrace() errors.StackTrace
			}
			errChain := []error{}
			for it := err; it != nil; it = errors.Unwrap(it) {
				errChain = append(errChain, it)
			}
			for index := len(errChain) - 1; index > 0; index-- {
				it := errChain[index]
				errWithStack, ok := it.(stackTracer)
				if ok {
					app.BaseCLI.Console().Warnf("Error stack trace:%+v\n", errWithStack.StackTrace())
					break
				}
			}
		}

		switch {
		case runExitCodeRegexp.MatchString(err.Error()):
			// This error would have been displayed earlier from the SolverMonitor.
			// This SetFatalError is a catch-all just in case that hasn't happened.
			app.BaseCLI.Logbus().Run().SetFatalError(
				time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_OTHER,
				err.Error())
			return 1
		case strings.Contains(err.Error(), "security.insecure is not allowed"):
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_NEEDS_PRIVILEGED, err.Error())
			app.BaseCLI.Console().Warnf("Error: earthly --allow-privileged (earthly -P) flag is required\n")
			return 9
		case strings.Contains(err.Error(), errutil.EarthlyGitStdErrMagicString):
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_GIT, err.Error())
			gitStdErr, shorterErr, ok := errutil.ExtractEarthlyGitStdErr(err.Error())
			if ok {
				app.BaseCLI.Console().Warnf("Error: %v\n\n%s\n", shorterErr, gitStdErr)
			} else {
				app.BaseCLI.Console().Warnf("Error: %v\n", err.Error())
			}
			app.BaseCLI.Console().Printf(
				"Check your git auth settings.\n" +
					"Did you ssh-add today? Need to configure ~/.earthly/config.yml?\n" +
					"For more information see https://docs.earthly.dev/guides/auth\n")
			return 1
		case strings.Contains(err.Error(), "failed to compute cache key") && strings.Contains(err.Error(), ": not found"):
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND, err.Error())
			re := regexp.MustCompile(`("[^"]*"): not found`)
			var matches = re.FindStringSubmatch(err.Error())
			if len(matches) == 2 {
				app.BaseCLI.Console().Warnf("Error: File not found %v\n", matches[1])
			} else {
				app.BaseCLI.Console().Warnf("Error: File not found: %v\n", err.Error())
			}
			return 1
		case strings.Contains(err.Error(), "429 Too Many Requests"):
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_RATE_LIMITED, err.Error())
			var registryName string
			var registryHost string
			if strings.Contains(err.Error(), "docker.com/increase-rate-limit") {
				registryName = "DockerHub"
			} else {
				registryName = "The remote registry"
				registryHost = " <server>" // keep the leading space
			}
			app.BaseCLI.Console().Warnf("Error: %s responded with a rate limit error. This is usually because you are not logged in.\n"+
				"You can login using the command:\n"+
				"  docker login%s", registryName, registryHost)
			return 1
		case strings.Contains(failedOutput, "Invalid ELF image for this architecture"):
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_OTHER, err.Error())
			app.BaseCLI.Console().Printf(
				"Are you using --platform to target a different architecture? You may have to manually install QEMU.\n" +
					"For more information see https://docs.earthly.dev/guides/multi-platform\n")
			return 1
		case !app.BaseCLI.Flags().Verbose && rpcRegex.MatchString(err.Error()):
			baseErr := errors.Cause(err)
			baseErrMsg := rpcRegex.ReplaceAllString(baseErr.Error(), "")
			app.BaseCLI.Console().Warnf("Error: %s\n", string(baseErrMsg))
			if strings.Contains(baseErrMsg, "transport is closing") {
				app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED, baseErr.Error())
				app.BaseCLI.Console().Warnf(
					"It seems that buildkitd is shutting down or it has crashed. " +
						"You can report crashes at https://github.com/earthly/earthly/issues/new.")
				if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) {
					app.printCrashLogs(ctx)
				}
				return 7
			} else {
				app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_OTHER, err.Error())
				return 1
			}
		case errors.Is(err, buildkitd.ErrBuildkitCrashed):
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED, err.Error())
			app.BaseCLI.Console().Warnf("Error: %v\n", err)
			app.BaseCLI.Console().Warnf(
				"It seems that buildkitd is shutting down or it has crashed. " +
					"You can report crashes at https://github.com/earthly/earthly/issues/new.")
			if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) {
				app.printCrashLogs(ctx)
			}
			return 7
		case errors.Is(err, buildkitd.ErrBuildkitConnectionFailure):
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_CONNECTION_FAILURE, err.Error())
			app.BaseCLI.Console().Warnf("Error: %v\n", err)
			if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) {
				app.BaseCLI.Console().Warnf(
					"It seems that buildkitd had an issue. " +
						"You can report crashes at https://github.com/earthly/earthly/issues/new.")
				app.printCrashLogs(ctx)
			}
			return 6
		case errors.Is(err, context.Canceled):
			app.BaseCLI.Logbus().Run().SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_CANCELED)
			app.BaseCLI.Console().Warnf("Canceled\n")
			app.BaseCLI.Console().VerbosePrintf("Canceled: %v\n", err)
			return 2
		case status.Code(errors.Cause(err)) == codes.Canceled:
			app.BaseCLI.Logbus().Run().SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_CANCELED)
			app.BaseCLI.Console().Warnf("Canceled\n")
			app.BaseCLI.Console().VerbosePrintf("Canceled: %v\n", err)
			if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) {
				app.printCrashLogs(ctx)
			}
			return 2
		case isInterpreterError:
			app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), ie.TargetID, "", logstream.FailureType_FAILURE_TYPE_SYNTAX, ie.Error())
			app.BaseCLI.Console().Warnf("Error: %s\n", ie.Error())
			return 1
		default:
			if app.BaseCLI.CommandName() == "build" {
				app.BaseCLI.Logbus().Run().SetFatalError(time.Now(), "", "", logstream.FailureType_FAILURE_TYPE_OTHER, err.Error())
			} else {
				app.BaseCLI.Logbus().Run().SkipFatalError()
			}
			app.BaseCLI.Console().Warnf("Error: %v\n", err)
			return 1
		}
	}
	app.BaseCLI.Logbus().Run().SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_SUCCESS)
	return 0
}

func (app *EarthlyApp) printCrashLogs(ctx context.Context) {
	app.BaseCLI.Console().PrintBar(color.New(color.FgHiRed), "System Info", "")
	fmt.Fprintf(os.Stderr, "version: %s\n", app.BaseCLI.Version())
	fmt.Fprintf(os.Stderr, "build-sha: %s\n", app.BaseCLI.GitSHA())
	fmt.Fprintf(os.Stderr, "platform: %s\n", common.GetPlatform())

	dockerVersion, err := buildkitd.GetDockerVersion(ctx, app.BaseCLI.Flags().ContainerFrontend)
	if err != nil {
		app.BaseCLI.Console().Warnf("failed querying docker version: %s\n", err.Error())
	} else {
		app.BaseCLI.Console().PrintBar(color.New(color.FgHiRed), "Docker Version", "")
		fmt.Fprintln(os.Stderr, dockerVersion)
	}

	logs, err := buildkitd.GetLogs(ctx, app.BaseCLI.Flags().ContainerName, app.BaseCLI.Flags().ContainerFrontend, app.BaseCLI.Flags().BuildkitdSettings)
	if err != nil {
		app.BaseCLI.Console().Warnf("failed fetching earthly-buildkit logs: %s\n", err.Error())
	} else {
		app.BaseCLI.Console().PrintBar(color.New(color.FgHiRed), "Buildkit Logs", "")
		fmt.Fprintln(os.Stderr, logs)
	}
}
