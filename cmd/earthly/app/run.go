package app

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/fatih/color"
	"github.com/moby/buildkit/util/grpcerrors"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/codes"

	"github.com/earthly/earthly/buildkitd"
	"github.com/earthly/earthly/cmd/earthly/common"
	"github.com/earthly/earthly/cmd/earthly/helper"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/earthfile2llb"
	"github.com/earthly/earthly/inputgraph"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/errutil"
	"github.com/earthly/earthly/util/hint"
	"github.com/earthly/earthly/util/params"
	"github.com/earthly/earthly/util/reflectutil"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/util/syncutil"
	"google.golang.org/grpc/status"
)

var (
	runExitCodeRegex  = regexp.MustCompile(`did not complete successfully: exit code: [^0][0-9]*($|[\n\t]+in\s+.*?\+.+)`)
	notFoundRegex     = regexp.MustCompile(`("[^"]*"): not found`)
	qemuExitCodeRegex = regexp.MustCompile(`process "/dev/.buildkit_qemu_emulator.*?did not complete successfully: exit code: 255$`)
	maxExecTimeRegex  = regexp.MustCompile(`max execution time of .+ exceeded`)
	requestIDRegex    = regexp.MustCompile(`(?P<msg>.*?) {reqID: .*?}`)
)

func (app *EarthlyApp) Run(ctx context.Context, console conslogging.ConsoleLogger, startTime time.Time, lastSignal *syncutil.Signal) int {
	err := app.unhideFlags(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error un-hiding flags %v", err)
		os.Exit(1)
	}
	helper.AutoComplete(ctx, app.BaseCLI)

	exitCode := app.run(ctx, os.Args, lastSignal)

	// app.Cfg will be nil when a user runs `earthly --version`;
	// however in all other regular commands app.Cfg will be set in app.Before
	if !app.BaseCLI.Flags().DisableAnalytics && app.BaseCLI.Cfg() != nil && !app.BaseCLI.Cfg().Global.DisableAnalytics {
		// Use a new context, in case the original context is cancelled due to sigint.
		_, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		displayErrors := app.BaseCLI.Flags().Verbose
		_, err := helper.NewCloudClient(app.BaseCLI)
		if err != nil && displayErrors {
			app.BaseCLI.Console().Warnf("unable to start cloud app.BaseClient: %s", err)
		}
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

func (app *EarthlyApp) run(ctx context.Context, args []string, lastSignal *syncutil.Signal) int {
	defer func() {
		if app.BaseCLI.LogbusSetup() != nil {
			err := app.BaseCLI.LogbusSetup().Close(ctx)
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
	defer app.BaseCLI.ExecuteDeferredFuncs()
	app.BaseCLI.Logbus().Run().SetStart(time.Now())

	defer func() {
		// Just in case this is forgotten somewhere else.
		app.BaseCLI.Logbus().Run().SetGenericFatalError(
			time.Now(),
			logstream.FailureType_FAILURE_TYPE_OTHER,
			"",
			"Error: No SetFatalError called appropriately. This should never happen.",
		)
	}()

	err := app.BaseCLI.App().RunContext(ctx, args)
	if err != nil {
		ie, isInterpreterError := earthfile2llb.GetInterpreterError(err)
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

		grpcErr, grpcErrOK := grpcerrors.AsGRPCStatus(err)
		hintErr, hintErrOK := getHintErr(err, grpcErr)
		var paramsErr *params.Error
		var autoSkipErr *inputgraph.Error
		switch {
		case hintErrOK:
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_OTHER,
				hintErr.Hint(),
				hintErr.Message(),
			)
			app.BaseCLI.Console().HelpPrintf(hintErr.Hint())
			return 1
		case errors.As(err, &autoSkipErr):
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_AUTO_SKIP,
				"",
				inputgraph.FormatError(err),
			)
			return 1
		case errors.As(err, &paramsErr):
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_INVALID_PARAM,
				"",
				paramsErr.ParentError(),
			)
			if paramsErr.Error() != paramsErr.ParentError() {
				app.BaseCLI.Console().VerboseWarnf(errorWithPrefix(paramsErr.Error()))
			}
			return 1
		case qemuExitCodeRegex.MatchString(err.Error()):
			var helpMsg string
			if app.BaseCLI.AnaMetaIsSat() {
				app.BaseCLI.Console().DebugPrintf("Are you using --platform to target a different architecture? Please note that \"disable-emulation\" flag is set in your satellite.\n")
			} else {
				helpMsg = "Are you using --platform to target a different architecture? You may have to manually install QEMU.\n" +
					"For more information see https://docs.earthly.dev/guides/multi-platform\n"
				app.BaseCLI.Console().HelpPrintf(helpMsg)
			}
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_OTHER,
				helpMsg,
				err.Error(),
			)
			return 255
		case runExitCodeRegex.MatchString(err.Error()):
			var helpMsg string
			if !app.BaseCLI.Flags().InteractiveDebugging && len(args) > 0 {
				args = append([]string{args[0], "-i"}, args[1:]...)
				args = redactSecretsFromArgs(args)
				args = stringutil.FilterElementsFromList(args, "--ci")
				msg := "To debug your build, you can use the --interactive (-i) flag to drop into a shell of the failing RUN step"
				helpMsg = fmt.Sprintf("%s: %q\n", msg, strings.Join(args, " "))
				app.BaseCLI.Console().HelpPrintf(helpMsg)
			}
			// This error would have been displayed earlier from the SolverMonitor.
			// This SetGenericFatalError is a catch-all just in case that hasn't happened.
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_OTHER,
				helpMsg,
				err.Error(),
			)
			return 1
		case strings.Contains(err.Error(), "security.insecure is not allowed"):
			helpMsg := "earthly --allow-privileged (earthly -P) flag is required\n"
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_NEEDS_PRIVILEGED,
				helpMsg,
				err.Error(),
			)
			app.BaseCLI.Console().HelpPrintf(helpMsg)
			return 9
		case strings.Contains(err.Error(), errutil.EarthlyGitStdErrMagicString):
			helpMsg := "Check your git auth settings.\n" +
				"Did you ssh-add today? Need to configure ~/.earthly/config.yml?\n" +
				"For more information see https://docs.earthly.dev/guides/auth\n"
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_GIT,
				helpMsg,
				err.Error(),
			)
			gitStdErr, shorterErr, ok := errutil.ExtractEarthlyGitStdErr(err.Error())
			if ok {
				app.BaseCLI.Console().VerboseWarnf("Error: %v\n\n%s\n", shorterErr, gitStdErr)
			} else {
				app.BaseCLI.Console().VerboseWarnf("Error: %v\n", err.Error())
			}
			app.BaseCLI.Console().HelpPrintf(helpMsg)
			return 1
		case strings.Contains(err.Error(), "failed to compute cache key") && strings.Contains(err.Error(), ": not found"):
			var matches = notFoundRegex.FindStringSubmatch(err.Error())
			msg := ""
			if len(matches) == 2 {
				msg = fmt.Sprintf("File not found: %s, %s\n", matches[1], err.Error())
			} else {
				msg = fmt.Sprintf("File not found: %s\n", err.Error())
			}
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND,
				"",
				msg,
			)
			return 1
		case strings.Contains(err.Error(), "429 Too Many Requests"):
			var registryName, registryHost string
			if strings.Contains(err.Error(), "docker.com/increase-rate-limit") {
				registryName = "DockerHub"
			} else {
				registryName = "The remote registry"
				registryHost = " <server>" // keep the leading space
			}
			helpMsg := fmt.Sprintf("%s responded with a rate limit error. This is usually because you are not logged in.\n"+
				"You can login using the command:\n"+
				"  docker login%s", registryName, registryHost)
			app.BaseCLI.Console().HelpPrintf(helpMsg)
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_RATE_LIMITED,
				helpMsg,
				err.Error(),
			)
			return 1
		case grpcErrOK && grpcErr.Code() == codes.PermissionDenied && requestIDRegex.MatchString(grpcErr.Message()):
			errorMsg := grpcErr.Message()
			matches, _ := stringutil.NamedGroupMatches(errorMsg, requestIDRegex)
			if len(matches["msg"]) > 0 {
				errorMsg = matches["msg"][0]
			}
			app.BaseCLI.Console().VerboseWarnf(err.Error())
			app.BaseCLI.Logbus().Run().SetGenericFatalError(time.Now(), logstream.FailureType_FAILURE_TYPE_OTHER, "", errorMsg)
			return 1
		case grpcErrOK && grpcErr.Code() == codes.Unknown && maxExecTimeRegex.MatchString(grpcErr.Message()):
			app.BaseCLI.Console().VerboseWarnf(errorWithPrefix(err.Error()))
			helpMsg := "Unverified accounts have a limit on the duration of RUN commands. Verify your account to lift this restriction."
			app.BaseCLI.Logbus().Run().SetGenericFatalError(time.Now(), logstream.FailureType_FAILURE_TYPE_OTHER, helpMsg, grpcErr.Message())
			app.BaseCLI.Console().HelpPrintf(helpMsg)
			return 1
		case grpcErrOK && grpcErr.Code() != codes.Canceled:
			app.BaseCLI.Console().VerboseWarnf(errorWithPrefix(err.Error()))
			if !strings.Contains(grpcErr.Message(), "transport is closing") {
				app.BaseCLI.Logbus().Run().SetGenericFatalError(
					time.Now(),
					logstream.FailureType_FAILURE_TYPE_OTHER,
					"",
					grpcErr.Message(),
				)
				return 1
			}
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED,
				"",
				grpcErr.Message(),
			)
			app.BaseCLI.Console().Warnf(
				"Error: It seems that buildkitd is shutting down or it has crashed. " +
					"You can report crashes at https://github.com/earthly/earthly/issues/new.")
			if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) {
				app.printCrashLogs(ctx)
			}
			return 7
		case errors.Is(err, buildkitd.ErrBuildkitCrashed):
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_BUILDKIT_CRASHED,
				"",
				err.Error(),
			)
			app.BaseCLI.Console().Warnf(
				"Error: It seems that buildkitd is shutting down or it has crashed. " +
					"You can report crashes at https://github.com/earthly/earthly/issues/new.")
			if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) {
				app.printCrashLogs(ctx)
			}
			return 7
		case errors.Is(err, buildkitd.ErrBuildkitConnectionFailure):
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_CONNECTION_FAILURE,
				"",
				err.Error(),
			)
			if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) {
				app.BaseCLI.Console().Warnf(
					"Error: It seems that buildkitd had an issue. " +
						"You can report crashes at https://github.com/earthly/earthly/issues/new.")
				app.printCrashLogs(ctx)
			}
			return 6
		case errors.Is(err, context.Canceled), grpcErrOK && grpcErr.Code() == codes.Canceled:
			app.BaseCLI.Logbus().Run().SetEnd(time.Now(), logstream.RunStatus_RUN_STATUS_CANCELED)
			if app.BaseCLI.Flags().Verbose {
				app.BaseCLI.Console().Warnf("Canceled: %v\n", err)
			} else {
				app.BaseCLI.Console().Warnf("Canceled\n")
			}
			if containerutil.IsLocal(app.BaseCLI.Flags().BuildkitdSettings.BuildkitAddress) && lastSignal.Get() == nil {
				app.printCrashLogs(ctx)
			}
			return 2
		case isInterpreterError:
			if ie.TargetID == "" {
				app.BaseCLI.Logbus().Run().SetGenericFatalError(
					time.Now(),
					logstream.FailureType_FAILURE_TYPE_SYNTAX,
					"",
					ie.Error(),
				)
				return 1
			}
			app.BaseCLI.Logbus().Run().SetFatalError(
				time.Now(),
				ie.TargetID,
				"",
				logstream.FailureType_FAILURE_TYPE_SYNTAX,
				"",
				ie.Error(),
			)
			return 1
		default:
			app.BaseCLI.Logbus().Run().SetGenericFatalError(
				time.Now(),
				logstream.FailureType_FAILURE_TYPE_OTHER,
				"",
				err.Error(),
			)
			return 1
		}
		//  if this happens there is a missing "return" in the above switch statement.
		panic("error was not caught by error-handling switch") // nolint:all
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

func errorWithPrefix(err string) string {
	return fmt.Sprintf("Error: %s", err)
}

func getHintErr(err error, grpcError *status.Status) (*hint.Error, bool) {
	if res := new(hint.Error); errors.As(err, &res) {
		return res, true
	}
	if grpcError != nil {
		return hint.FromError(errors.New(grpcError.Message()))
	}
	return nil, false
}

func redactSecretsFromArgs(args []string) []string {
	redacted := []string{}
	isSecret := false
	for _, arg := range args {
		if isSecret {
			isSecret = false
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) > 1 {
				redacted = append(redacted, fmt.Sprintf("%s=XXXXX", parts[0]))
				continue
			}
		}
		if arg == "-s" || arg == "--secret" {
			isSecret = true
		}
		redacted = append(redacted, arg)
	}
	return redacted
}
