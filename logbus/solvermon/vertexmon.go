package solvermon

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/util/errutil"
	"github.com/earthly/earthly/util/statsstreamparser"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/util/vertexmeta"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
)

const (
	// BuildkitStatsStream is the stream number associated with runc stats
	BuildkitStatsStream = 99 // TODO move to a common location in buildkit
)

type vertexMonitor struct {
	vertex    *client.Vertex
	meta      *vertexmeta.VertexMeta
	operation string
	cp        *logbus.Command
	ssp       *statsstreamparser.Parser

	isFatalError   bool // If set, this is the root cause of the entire build failure.
	fatalErrorType logstream.FailureType
	errorStr       string
	isCanceled     bool
}

var reErrExitCode = regexp.MustCompile(`^(?:process ".*" did not complete successfully|error calling LocalhostExec): exit code: (?P<exit_code>[0-9]+)$`)

func getExitCode(errString string) (uint64, bool) {
	if reErrExitCode.MatchString(errString) {
		matches, _ := stringutil.NamedGroupMatches(errString, reErrExitCode)
		exitCodeMatch := matches["exit_code"][0]
		exitCode, err := strconv.ParseUint(exitCodeMatch, 10, 32)
		if err != nil {
			return 0, false
		}
		return exitCode, true
	}
	return 0, false
}

var reErrNotFound = regexp.MustCompile(`^failed to calculate checksum of ref ([^ ]*): (.*)$`)
var reHint = regexp.MustCompile(`^(?P<msg>.+?):Hint: .+`)

func determineFatalErrorType(errString string, exitCode uint64) (logstream.FailureType, bool) {
	if strings.Contains(errString, "context canceled") || errString == "no active sessions" {
		return 0, false
	}
	if reErrExitCode.MatchString(errString) {
		switch exitCode {
		case math.MaxUint32:
			return logstream.FailureType_FAILURE_TYPE_OOM_KILLED, true
		default:
			return logstream.FailureType_FAILURE_TYPE_NONZERO_EXIT, true
		}
	}
	if reErrNotFound.MatchString(errString) {
		return logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND, true
	}
	if strings.Contains(errString, errutil.EarthlyGitStdErrMagicString) {
		return logstream.FailureType_FAILURE_TYPE_GIT, true
	}
	return logstream.FailureType(0), false
}

func formatError(errString string, fatalErrorType logstream.FailureType, exitCode uint64, indentOp string, internalStr string) string {
	if matches, _ := stringutil.NamedGroupMatches(errString, reHint); len(matches["msg"]) == 1 {
		errString = matches["msg"][0]
	}
	formattedError := ""
	switch fatalErrorType {
	case logstream.FailureType_FAILURE_TYPE_OOM_KILLED:
		formattedError = fmt.Sprintf(""+
			"      The%s command\n"+
			"          %s\n"+
			"      was terminated because the build system ran out of memory. If you are using a satellite or other remote buildkit, it is the remote system that ran out of memory.",
			internalStr, indentOp)
	case logstream.FailureType_FAILURE_TYPE_NONZERO_EXIT:
		formattedError = fmt.Sprintf(""+
			"      The%s command\n"+
			"          %s\n"+
			"      did not complete successfully. Exit code %d",
			internalStr, indentOp, exitCode)
	case logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND:
		m := reErrNotFound.FindStringSubmatch(errString)
		formattedError = fmt.Sprintf(""+
			"      The%s command\n"+
			"          %s\n"+
			"      failed: %s",
			internalStr, indentOp, m[2])
	case logstream.FailureType_FAILURE_TYPE_GIT:
		gitStdErr, shorterErr, ok := errutil.ExtractEarthlyGitStdErr(errString)
		if ok {
			formattedError = fmt.Sprintf("The%s command\n          %s\nfailed: %s\n\n%s", internalStr, indentOp, shorterErr, gitStdErr)
		} else {
			formattedError = fmt.Sprintf("The%s command\n          %s\nfailed: %s", internalStr, indentOp, errString)
		}
	default:
		formattedError = fmt.Sprintf("The%s command\n          %s\nfailed: %s", internalStr, indentOp, errString)
	}
	return formattedError
}

func FormatError(errString string) string {
	exitCode, _ := getExitCode(errString)
	fatalErrorType, _ := determineFatalErrorType(errString, exitCode)
	return formatError(errString, fatalErrorType, exitCode, "", "")
}

func (vm *vertexMonitor) parseError() {
	errString := vm.vertex.Error

	// Add operation context to the error string
	indentOp := strings.Join(strings.Split(vm.operation, "\n"), "\n          ")
	internalStr := ""
	if vm.meta.Internal {
		internalStr = " internal"
	}
	errString = fmt.Sprintf("%s%s", internalStr, errString)

	exitCode, _ := getExitCode(errString)
	fatalErrorType, isFatalError := determineFatalErrorType(errString, exitCode)
	formattedError := formatError(errString, fatalErrorType, exitCode, indentOp, internalStr)

	// Add source location if available
	slString := ""
	if vm.meta.SourceLocation != nil {
		slString = fmt.Sprintf(
			" %s:%d:%d",
			vm.meta.SourceLocation.File, vm.meta.SourceLocation.StartLine,
			vm.meta.SourceLocation.StartColumn)
	}

	// Set the error string and flags on the vertexMonitor
	if isFatalError {
		vm.errorStr = fmt.Sprintf("ERROR%s\n%s", slString, formattedError)
	} else {
		vm.errorStr = fmt.Sprintf("WARN%s: %s", slString, formattedError)
	}
	vm.isFatalError = isFatalError
	vm.fatalErrorType = fatalErrorType
}

func (vm *vertexMonitor) Write(dt []byte, ts time.Time, stream int) (int, error) {
	if stream == BuildkitStatsStream {
		stats, err := vm.ssp.Parse(dt)
		if err != nil {
			return 0, errors.Wrap(err, "failed decoding stats stream")
		}
		for _, statsSample := range stats {
			statsJSON, err := json.Marshal(statsSample)
			if err != nil {
				return 0, errors.Wrap(err, "stats json encode failed")
			}
			_, err = vm.cp.Write(statsJSON, ts, int32(stream))
			if err != nil {
				return 0, errors.Wrap(err, "write stats")
			}
		}
		return len(dt), nil
	}
	_, err := vm.cp.Write(dt, ts, int32(stream))
	if err != nil {
		return 0, errors.Wrap(err, "write log line")
	}
	return len(dt), nil
}
