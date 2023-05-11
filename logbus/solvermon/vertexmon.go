package solvermon

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/util/errutil"
	"github.com/earthly/earthly/util/vertexmeta"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
)

type vertexMonitor struct {
	vertex    *client.Vertex
	meta      *vertexmeta.VertexMeta
	operation string
	cp        *logbus.Command

	isFatalError   bool // If set, this is the root cause of the entire build failure.
	fatalErrorType logstream.FailureType
	errorStr       string
	isCanceled     bool
}

var reErrExitCode = regexp.MustCompile(`^process (".*") did not complete successfully: exit code: ([0-9]+)$`)
var reErrNotFound = regexp.MustCompile(`^failed to calculate checksum of ref ([^ ]*): (.*)$`)

func (vm *vertexMonitor) Write(dt []byte, ts time.Time, stream int) (int, error) {
	_, err := vm.cp.Write(dt, ts, int32(stream))
	if err != nil {
		return 0, errors.Wrap(err, "write log line")
	}
	return len(dt), nil
}

func (vm *vertexMonitor) parseError() {
	errString := vm.vertex.Error
	indentOp := strings.Join(strings.Split(vm.operation, "\n"), "\n          ")
	internalStr := ""
	if vm.meta.Internal {
		internalStr = " internal"
	}
	switch {
	case strings.Contains(errString, "context canceled"):
		vm.isCanceled = true
		vm.errorStr = "WARN: Canceled"
		return
	case reErrExitCode.MatchString(errString):
		m := reErrExitCode.FindStringSubmatch(errString)

		// Ignore the parse error as default case will print it as a string using
		// the source, so we won't miss any data.
		exitCode, _ := strconv.ParseUint(m[2], 10, 32)
		switch exitCode {
		case math.MaxUint32:
			errString = fmt.Sprintf(""+
				"      The%s command\n"+
				"          %s\n"+
				"      was terminated because the build system ran out of memory.\n"+
				"      If you are using a satellite or other remote buildkit, it is the remote system that ran out of memory.",
				internalStr, indentOp)
			vm.fatalErrorType = logstream.FailureType_FAILURE_TYPE_OOM_KILLED
		default:
			errString = fmt.Sprintf(""+
				"      The%s command\n"+
				"          %s\n"+
				"      did not complete successfully. Exit code %s",
				internalStr, indentOp, m[2])
			vm.fatalErrorType = logstream.FailureType_FAILURE_TYPE_NONZERO_EXIT
		}
		vm.isFatalError = true
	case reErrNotFound.MatchString(errString):
		m := reErrNotFound.FindStringSubmatch(errString)
		errString = fmt.Sprintf(""+
			"      The%s command\n"+
			"          %s\n"+
			"      failed: %s",
			internalStr, indentOp, m[2])
		vm.isFatalError = true
		vm.fatalErrorType = logstream.FailureType_FAILURE_TYPE_FILE_NOT_FOUND
	case errString == "no active sessions":
		vm.isCanceled = true
		errString = "WARN: Canceled"
	case strings.Contains(errString, errutil.EarthlyGitStdErrMagicString):
		gitStdErr, shorterErr, ok := errutil.ExtractEarthlyGitStdErr(errString)
		if ok {
			errString = fmt.Sprintf(
				"The%s command '%s' failed: %s\n\n%s", internalStr, vm.operation, shorterErr, gitStdErr)
		} else {
			errString = fmt.Sprintf(
				"The%s command '%s' failed: %s", internalStr, vm.operation, errString)
		}
		vm.isFatalError = true
	default:
		errString = fmt.Sprintf(
			"The%s command '%s' failed: %s", internalStr, vm.operation, errString)
	}
	slString := ""
	if vm.meta.SourceLocation != nil {
		slString = fmt.Sprintf(
			" %s line %d:%d",
			vm.meta.SourceLocation.File, vm.meta.SourceLocation.StartLine,
			vm.meta.SourceLocation.StartColumn)
	}
	if vm.isFatalError {
		vm.errorStr = fmt.Sprintf("ERROR%s\n%s", slString, errString)
	} else {
		vm.errorStr = fmt.Sprintf("WARN%s: %s", slString, errString)
	}
}
