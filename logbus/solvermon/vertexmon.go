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

var reErrExitCode = regexp.MustCompile(`(?:process ".*" did not complete successfully|error calling LocalhostExec): exit code: (?P<exit_code>[0-9]+)$`)

var errNoExitCodeOMM = errors.New("no exit code, process was killed due to OOM")
var errNoExitCode = errors.New("no exit code in error message")

// getExitCode returns the exit code (int), whether one was found (bool), and an error if the exit code was invalid
func getExitCode(errString string) (int, error) {
	if matches, _ := stringutil.NamedGroupMatches(errString, reErrExitCode); len(matches["exit_code"]) == 1 {
		exitCodeMatch := matches["exit_code"][0]
		exitCode, err := strconv.ParseUint(exitCodeMatch, 10, 32)
		if err != nil {
			return 0, err
		}
		// See https://github.com/earthly/buildkit/commit/9b0bdb600641f3dd1d96f54ac2d86581ab6433b2
		if exitCode == math.MaxUint32 {
			return 0, errNoExitCodeOMM
		}
		if exitCode > 255 {
			return 0, fmt.Errorf("exit code %d out of expected range (0-255)", exitCode)
		}
		exitCodeByte := exitCode & 0xFF
		return int(exitCodeByte), nil
	}
	return 0, errNoExitCode
}

var reErrNotFound = regexp.MustCompile(`^\s*(internal)?failed to calculate checksum of ref ([^ ]::[^ ]*|[^ ]*): (.*)\s*$`)
var reHint = regexp.MustCompile(`^(?P<msg>.+?):Hint: .+`)

// determineFatalErrorType returns logstream.FailureType
// and whether or not its a Fatal Error
func determineFatalErrorType(errString string, exitCode int, exitParseErr error) (logstream.FailureType, bool) {
	if strings.Contains(errString, "context canceled") || errString == "no active sessions" {
		return logstream.FailureType_FAILURE_TYPE_UNKNOWN, false
	}
	if exitParseErr == errNoExitCodeOMM {
		return logstream.FailureType_FAILURE_TYPE_OOM_KILLED, true
	} else if exitParseErr != nil && exitParseErr != errNoExitCode {
		// We have an exit code, and can't parse it
		return logstream.FailureType_FAILURE_TYPE_UNKNOWN, true
	}
	if exitCode > 0 {
		return logstream.FailureType_FAILURE_TYPE_NONZERO_EXIT, true
	}
	if reErrNotFound.MatchString(errString) {
		return logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND, true
	}
	if strings.Contains(errString, errutil.EarthlyGitStdErrMagicString) {
		return logstream.FailureType_FAILURE_TYPE_GIT, true
	}
	return logstream.FailureType_FAILURE_TYPE_UNKNOWN, false
}

func formatErrorMessage(errString, operation string, internal bool, fatalErrorType logstream.FailureType, exitCode int) string {
	if matches, _ := stringutil.NamedGroupMatches(errString, reHint); len(matches["msg"]) == 1 {
		errString = matches["msg"][0]
	}

	internalStr := ""
	if internal {
		internalStr = " internal"
	}
	errString = fmt.Sprintf("%s%s", internalStr, errString)

	switch fatalErrorType {
	case logstream.FailureType_FAILURE_TYPE_OOM_KILLED:
		return fmt.Sprintf(
			"      The%s command\n"+
				"          %s\n"+
				"      was terminated because the build system ran out of memory. If you are using a satellite or other remote buildkit, it is the remote system that ran out of memory.", internalStr, operation)
	case logstream.FailureType_FAILURE_TYPE_NONZERO_EXIT:
		return fmt.Sprintf(
			"      The%s command\n"+
				"          %s\n"+
				"      did not complete successfully. Exit code %d", internalStr, operation, exitCode)
	case logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND:
		m := reErrNotFound.FindStringSubmatch(errString)
		reason := fmt.Sprintf("unable to parse file_not_found error:%s", errString)
		if len(m) > 2 {
			reason = m[3]
		}
		return fmt.Sprintf(
			"      The%s command\n"+
				"          %s\n"+
				"      failed: %s", internalStr, operation, reason)
	case logstream.FailureType_FAILURE_TYPE_GIT:
		gitStdErr, shorterErr, ok := errutil.ExtractEarthlyGitStdErr(errString)
		if ok {
			return fmt.Sprintf(
				"The%s command\n"+
					"          %s\n"+
					"failed: %s\n\n%s", internalStr, operation, shorterErr, gitStdErr)
		}
		return fmt.Sprintf(
			"The%s command\n"+
				"          %s\n"+
				"failed: %s", internalStr, operation, errString)
	default:
		return fmt.Sprintf(
			"The%s command\n"+
				"          %s\n"+
				"failed: %s", internalStr, operation, errString)
	}
}

func FormatError(operation string, errString string) string {
	exitCode, err := getExitCode(errString)
	fatalErrorType, _ := determineFatalErrorType(errString, exitCode, err)
	return formatErrorMessage(errString, operation, false, fatalErrorType, exitCode)
}

func (vm *vertexMonitor) parseError() {
	errString := vm.vertex.Error

	indentOp := strings.Join(strings.Split(vm.operation, "\n"), "\n          ")

	exitCode, err := getExitCode(errString)
	fatalErrorType, isFatalError := determineFatalErrorType(errString, exitCode, err)
	formattedError := formatErrorMessage(errString, indentOp, vm.meta.Internal, fatalErrorType, exitCode)

	// Add Error location
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
