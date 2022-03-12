package outmon

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"github.com/earthly/earthly/conslogging"
	"github.com/mattn/go-isatty"
	"github.com/moby/buildkit/client"
	"github.com/pkg/errors"
)

type vertexMonitor struct {
	vertex         *client.Vertex
	meta           *VertexMeta
	operation      string
	lastProgress   map[string]time.Time
	lastPercentage map[string]int
	console        conslogging.ConsoleLogger
	headerPrinted  bool
	isError        bool
	isCanceled     bool
	tailOutput     *circbuf.Buffer
	// Line of output that has not yet been terminated with a \n.
	openLine            []byte
	lastOpenLineUpdate  time.Time
	lastOpenLineSkipped bool
}

func (vm *vertexMonitor) printHeader() {
	vm.headerPrinted = true
	if vm.operation == "" {
		return
	}
	c := vm.console
	var metaParts []string
	if vm.meta.NonDefaultPlatform && vm.meta.Platform != "" {
		metaParts = append(metaParts, vm.meta.Platform)
	}
	if vm.meta.OverridingArgs != nil {
		metaParts = append(metaParts, vm.meta.OverridingArgsString())
	}
	if len(metaParts) > 0 {
		c.WithMetadataMode(true).Printf("%s\n", strings.Join(metaParts, " | "))
	}
	out := []string{}
	out = append(out, "-->")
	out = append(out, vm.operation)
	if vm.vertex.Cached {
		c = c.WithCached(true)
	}
	c.Printf("%s\n", strings.Join(out, " "))
}

var internalProgress = map[string]bool{
	"exporting manifest": true,
	"sending tarballs":   true,
	"exporting config":   true,
	"exporting layers":   true,
	"copying files":      true,
}

const esc = 27

var ansiUp = []byte(fmt.Sprintf("%c[A", esc))
var ansiEraseRestLine = []byte(fmt.Sprintf("%c[K", esc))
var ansiSupported = os.Getenv("TERM") != "dumb" &&
	(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()))

func (vm *vertexMonitor) printOutput(output []byte, sameAsLast bool) error {
	if vm.tailOutput == nil {
		var err error
		vm.tailOutput, err = circbuf.NewBuffer(tailErrorBufferSizeBytes)
		if err != nil {
			return errors.Wrap(err, "allocate buffer for output")
		}
	}
	// Use the raw output for the tail buffer.
	_, err := vm.tailOutput.Write(output)
	if err != nil {
		return errors.Wrap(err, "write to in-memory output buffer")
	}
	printOutput := make([]byte, 0, len(vm.openLine)+len(output)+10)
	if bytes.HasPrefix(output, []byte{'\n'}) && len(vm.openLine) > 0 && !vm.lastOpenLineSkipped {
		// Optimization for cases where ansi control sequences are not supported:
		// if the output starts with a \n, then treat the open line as closed and
		// just keep going after that.
		vm.openLine = nil
		output = output[1:]
		vm.lastOpenLineUpdate = time.Time{}
	}
	if sameAsLast && len(vm.openLine) > 0 {
		// Prettiness optimization: if there is an open line and the previous print out
		// was of the same vertex, then use ANSI control sequence to go up one line and
		// keep writing there.
		printOutput = append(printOutput, ansiUp...)
	}
	// Prepend the open line to the output.
	printOutput = append(printOutput, vm.openLine...)
	printOutput = append(printOutput, output...)
	// Look for the last \n to update the open line.
	lastNewLine := bytes.LastIndexByte(printOutput, '\n')
	if lastNewLine != -1 {
		// Ends up being empty slice if output ends in \n.
		vm.openLine = printOutput[(lastNewLine + 1):]
		// A \n exists - reset the open line timer.
		vm.lastOpenLineUpdate = time.Time{}
	} else {
		// No \n found - update vm.openLine to append the new output.
		vm.openLine = append(vm.openLine, output...)
	}
	if !bytes.HasSuffix(printOutput, []byte{'\n'}) {
		if vm.lastOpenLineUpdate.Add(durationBetweenOpenLineUpdate).After(time.Now()) {
			// Skip printing if trying to update the same line too frequently.
			vm.lastOpenLineSkipped = true
			return nil
		}
		vm.lastOpenLineUpdate = time.Now()
		// If output doesn't terminate in \n, add our own.
		printOutput = append(printOutput, '\n')
	}
	vm.lastOpenLineSkipped = false
	vm.console.PrintBytes(printOutput)
	return nil
}

func (vm *vertexMonitor) shouldPrintProgress(id string, percent int, verbose bool, sameAsLast bool) bool {
	if !vm.headerPrinted {
		return false
	}
	if !verbose && !ansiSupported {
		for prefix := range internalProgress {
			if strings.HasPrefix(id, prefix) {
				return false
			}
		}
	}
	minDelta := durationBetweenProgressUpdate
	if sameAsLast && ansiSupported {
		minDelta = durationBetweenProgressUpdateIfSame
	} else if strings.HasPrefix(id, "sha256:") || strings.HasPrefix(id, "extracting sha256:") {
		// These progress updates are a bit more annoying - do them more rarely.
		minDelta = durationBetweenSha256ProgressUpdate
	}
	now := time.Now()
	lastProgress := vm.lastProgress[id]
	lastPercentage := -1
	lastPercentageStored, ok := vm.lastPercentage[id]
	if ok {
		lastPercentage = lastPercentageStored
	}
	if now.Sub(lastProgress) < minDelta && percent < 100 {
		return false
	}
	if lastPercentage == percent {
		return false
	}
	vm.lastProgress[id] = now
	vm.lastPercentage[id] = percent
	return true
}

func (vm *vertexMonitor) printProgress(id string, progress int, verbose bool, sameAsLast bool) {
	builder := make([]string, 0, 2)
	if sameAsLast {
		// Overwrite previous line if this update is for the same thing as the previous one.
		builder = append(builder, string(ansiUp))
	}
	progressBar := progressBar(progress, 10)
	builder = append(builder, fmt.Sprintf("[%s] %3d%% %s%s\n", progressBar, progress, id, string(ansiEraseRestLine)))
	vm.console.PrintBytes([]byte(strings.Join(builder, "")))
}

func (vm *vertexMonitor) printError() bool {
	if strings.Contains(vm.vertex.Error, "did not complete successfully") {
		vm.console.Warnf("ERROR: Command exited with non-zero code: %s\n", vm.operation)
		return true
	}
	vm.console.Printf("WARN: (%s) %s\n", vm.operation, vm.vertex.Error)
	return false
}

func (vm *vertexMonitor) printTimingInfo() {
	if vm.vertex.Started == nil || vm.vertex.Completed == nil {
		return
	}
	vm.console.WithMetadataMode(true).
		Printf("Completed in %s\n", vm.vertex.Completed.Sub(*vm.vertex.Started))
}

func (vm *vertexMonitor) isOngoing() bool {
	return vm.vertex.Started != nil && vm.vertex.Completed == nil && !vm.isError
}

func (vm *vertexMonitor) reportStatusToConsole() {
	vm.console.MarkBundleBuilderStatus(vm.vertex.Started != nil, vm.vertex.Completed != nil, vm.isCanceled)
}

func (vm *vertexMonitor) reportResultToConsole() {
	vm.console.MarkBundleBuilderResult(vm.isError, vm.isCanceled)
}
