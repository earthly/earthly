package bus

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/armon/circbuf"
	"github.com/earthly/earthly/outmon"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
)

type vertexMonitor struct {
	vertex    *client.Vertex
	meta      *outmon.VertexMeta
	operation string
	cp        *CommandPrinter

	tailOutput *circbuf.Buffer

	isFatalError bool // If set, this is the root cause of the entire build failure.
	errorStr     string
	isCanceled   bool
}

var reErrExitCode = regexp.MustCompile(`^process (".*") did not complete successfully: exit code: ([0-9]+)$`)
var reErrNotFound = regexp.MustCompile(`^failed to calculate checksum of ref ([^ ]*): (.*)$`)

func (vm *vertexMonitor) Write(dt []byte) (int, error) {
	_, err := vm.cp.Write(dt)
	if err != nil {
		return 0, errors.Wrap(err, "write log line")
	}
	if vm.tailOutput == nil {
		var err error
		vm.tailOutput, err = circbuf.NewBuffer(tailErrorBufferSizeBytes)
		if err != nil {
			return 0, errors.Wrap(err, "allocate buffer for output")
		}
	}
	_, err = vm.tailOutput.Write(dt)
	if err != nil {
		return 0, errors.Wrap(err, "write to tail output buffer")
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
		errString = fmt.Sprintf(""+
			"      The%s command\n"+
			"          %s\n"+
			"      did not complete successfully. Exit code %s",
			internalStr, indentOp, m[2])
		vm.isFatalError = true
	case reErrNotFound.MatchString(errString):
		m := reErrNotFound.FindStringSubmatch(errString)
		errString = fmt.Sprintf(""+
			"      The%s command\n"+
			"          %s\n"+
			"      failed: %s",
			internalStr, indentOp, m[2])
		vm.isFatalError = true
	case errString == "no active sessions":
		errString = "Canceled"
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
