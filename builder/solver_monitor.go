package builder

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"github.com/earthly/earthly/conslogging"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

const (
	durationBetweenProgressUpdate = time.Second * 5
	durationBetweenOpenLineUpdate = time.Second
	tailErrorBufferSizeBytes      = 80 * 1024 // About as much as 1024 lines of 80 chars each.
)

type vertexMonitor struct {
	vertex         *client.Vertex
	targetStr      string
	targetBrackets string
	salt           string
	operation      string
	lastOutput     time.Time
	lastPercentage int
	console        conslogging.ConsoleLogger
	headerPrinted  bool
	isInternal     bool
	isError        bool
	tailOutput     *circbuf.Buffer
	// Line of output that has not yet been terminated with a \n.
	openLine            []byte
	lastOpenLineUpdate  time.Time
	lastOpenLineSkipped bool
}

func (vm *vertexMonitor) printHeader(printMetadata bool) {
	vm.headerPrinted = true
	if vm.operation == "" {
		return
	}
	c := vm.console
	if vm.targetBrackets != "" && printMetadata {
		c.WithMetadataMode(true).Printf("%s\n", vm.targetBrackets)
	}
	out := []string{}
	out = append(out, "-->")
	out = append(out, vm.operation)
	if vm.vertex.Cached {
		c = c.WithCached(true)
	}
	c.Printf("%s\n", strings.Join(out, " "))
}

func (vm *vertexMonitor) shouldPrintProgress(percent int) bool {
	if !vm.headerPrinted {
		return false
	}
	if vm.targetStr == "" {
		return false
	}
	now := time.Now()
	if now.Sub(vm.lastOutput) < durationBetweenProgressUpdate && percent < 100 {
		return false
	}
	if vm.lastPercentage >= percent {
		return false
	}
	vm.lastOutput = now
	vm.lastPercentage = percent
	return true
}

const esc = 27

var ansiUp = []byte(fmt.Sprintf("%c[A", esc))

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

func (vm *vertexMonitor) printError() {
	if strings.Contains(vm.vertex.Error, "executor failed running") {
		vm.console.Warnf("ERROR: Command exited with non-zero code: %s\n", vm.operation)
	} else {
		vm.console.Warnf("ERROR: (%s) %s\n", vm.operation, vm.vertex.Error)
	}
}

func (vm *vertexMonitor) printTimingInfo() {
	if vm.vertex.Started == nil || vm.vertex.Completed == nil {
		return
	}
	vm.console.WithMetadataMode(true).
		Printf("Completed in %s\n", vm.vertex.Completed.Sub(*vm.vertex.Started))
}

type solverMonitor struct {
	console          conslogging.ConsoleLogger
	verbose          bool
	vertices         map[digest.Digest]*vertexMonitor
	saltSeen         map[string]bool
	lastVertexOutput *vertexMonitor
	// timingTable is a map of target+salt string -> total duration.
	timingTable map[timingKey]time.Duration
	startTime   time.Time
}

type timingKey struct {
	targetStr      string
	targetBrackets string
	salt           string
}

func newSolverMonitor(console conslogging.ConsoleLogger, verbose bool) *solverMonitor {
	return &solverMonitor{
		console:     console,
		verbose:     verbose,
		vertices:    make(map[digest.Digest]*vertexMonitor),
		saltSeen:    make(map[string]bool),
		timingTable: make(map[timingKey]time.Duration),
		startTime:   time.Now(),
	}
}

func (sm *solverMonitor) monitorProgress(ctx context.Context, ch chan *client.SolveStatus) error {
	var errVertex *vertexMonitor
Loop:
	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				break Loop
			}
			for _, vertex := range ss.Vertexes {
				vm, ok := sm.vertices[vertex.Digest]
				if !ok {
					targetStr, targetBrackets, salt, operation := parseVertexName(vertex.Name)
					vm = &vertexMonitor{
						vertex:         vertex,
						targetStr:      targetStr,
						targetBrackets: targetBrackets,
						salt:           salt,
						operation:      operation,
						isInternal:     (targetStr == "internal" && !sm.verbose),
						console:        sm.console.WithPrefixAndSalt(targetStr, salt),
					}
					sm.vertices[vertex.Digest] = vm
				}
				vm.vertex = vertex
				if !vm.headerPrinted &&
					((!vm.isInternal && (vertex.Cached || vertex.Started != nil)) || vertex.Error != "") {
					sm.printHeader(vm)
				}
				if vertex.Error != "" {
					if strings.Contains(vertex.Error, "context canceled") {
						if !vm.isInternal {
							vm.console.Printf("WARN: Canceled\n")
						}
					} else {
						vm.isError = true
						if errVertex == nil {
							errVertex = vm
						}
						vm.printError()
					}
				}
				if sm.verbose {
					vm.printTimingInfo()
					sm.recordTiming(vm.targetStr, vm.targetBrackets, vm.salt, vertex)
				}
			}
			for _, vs := range ss.Statuses {
				vm, ok := sm.vertices[vs.Vertex]
				if !ok || vm.isInternal {
					// No logging for internal operations.
					continue
				}
				progress := int(0)
				if vs.Total != 0 {
					progress = int(100.0 * float32(vs.Current) / float32(vs.Total))
				}
				if vs.Completed != nil {
					progress = 100
				}
				if vm.shouldPrintProgress(progress) {
					if !vm.headerPrinted {
						sm.printHeader(vm)
					}
					vm.console.Printf("%s %d%%\n", vs.ID, progress)
				}
			}
			for _, logLine := range ss.Logs {
				vm, ok := sm.vertices[logLine.Vertex]
				if !ok || vm.isInternal {
					// No logging for internal operations.
					continue
				}
				if !vm.headerPrinted {
					sm.printHeader(vm)
				}
				err := sm.printOutput(vm, logLine.Data)
				if err != nil {
					return err
				}
			}
		}
	}
	if errVertex != nil {
		sm.reprintFailure(errVertex)
	}
	return nil
}

func (sm *solverMonitor) printOutput(vm *vertexMonitor, data []byte) error {
	sameAsLast := (sm.lastVertexOutput == vm)
	sm.lastVertexOutput = vm
	return vm.printOutput(data, sameAsLast)
}

func (sm *solverMonitor) printHeader(vm *vertexMonitor) {
	seen := sm.saltSeen[vm.salt]
	if !seen {
		sm.saltSeen[vm.salt] = true
	}
	vm.printHeader(!seen || sm.verbose)
}

func (sm *solverMonitor) recordTiming(targetStr, targetBrackets, salt string, vertex *client.Vertex) {
	if vertex.Started == nil || vertex.Completed == nil {
		return
	}
	dur := vertex.Completed.Sub(*vertex.Started)
	if dur == 0 {
		return
	}
	key := timingKey{
		targetStr:      targetStr,
		targetBrackets: targetBrackets,
		salt:           salt,
	}
	sm.timingTable[key] += dur
}

func (sm *solverMonitor) PrintTiming() {
	if !sm.verbose {
		return
	}
	sm.console.
		WithMetadataMode(true).
		Printf("Summary of timing information\n" +
			"Note that the times do not include the expansion of commands like BUILD, FROM, COPY (artifact).")
	var total time.Duration
	durs := make([]time.Duration, 0, len(sm.timingTable))
	durMap := make(map[time.Duration][]timingKey)
	for key, dur := range sm.timingTable {
		durs = append(durs, dur)
		durMap[dur] = append(durMap[dur], key)
		total += dur
	}
	sort.Slice(durs, func(i, j int) bool {
		return durs[i] > durs[j]
	})
	for _, dur := range durs {
		for _, key := range durMap[dur] {
			sm.console.
				WithPrefixAndSalt(key.targetStr, key.salt).
				WithMetadataMode(true).
				Printf("(%s) %s\n", key.targetBrackets, dur)
		}
	}
	sm.console.
		WithMetadataMode(true).
		Printf("===============================================================\n")
	sm.console.
		WithMetadataMode(true).
		Printf("Total       \t%s\n", total)
	sm.console.
		WithMetadataMode(true).
		Printf("Total (real)\t%s\n", time.Now().Sub(sm.startTime))
}

func (sm *solverMonitor) reprintFailure(errVertex *vertexMonitor) {
	sm.console.Warnf("Repeating the output of the command that caused the failure\n")
	sm.console.PrintFailure()
	errVertex.console = errVertex.console.WithFailed(true)
	errVertex.printHeader(true)
	if errVertex.tailOutput != nil {
		isTruncated := (errVertex.tailOutput.TotalWritten() > errVertex.tailOutput.Size())
		if errVertex.tailOutput.TotalWritten() == 0 {
			errVertex.console.Printf("[no output]\n")
		} else {
			if isTruncated {
				errVertex.console.Printf("[...]\n")
			}
			errVertex.console.PrintBytes(errVertex.tailOutput.Bytes())
		}
	} else {
		errVertex.console.Printf("[no output]\n")
	}
	errVertex.printError()
}

var vertexRegexp = regexp.MustCompile("^\\[([^\\]]*)\\] (.*)$")
var targetAndSaltRegexp = regexp.MustCompile("^([^\\(]*)(\\(([^\\)]*)\\))? (.*)$")

func parseVertexName(vertexName string) (string, string, string, string) {
	target := ""
	targetBrackets := ""
	operation := ""
	salt := ""
	match := vertexRegexp.FindStringSubmatch(vertexName)
	if len(match) < 2 {
		return "internal", targetBrackets, "internal", vertexName
	}
	targetAndSalt := match[1]
	operation = match[2]
	targetAndSaltMatch := targetAndSaltRegexp.FindStringSubmatch(targetAndSalt)
	if targetAndSaltMatch == nil {
		return targetAndSalt, targetBrackets, targetAndSalt, operation
	}
	target = targetAndSaltMatch[1]
	salt = targetAndSaltMatch[len(targetAndSaltMatch)-1]
	if salt == "" {
		salt = targetAndSalt
	}
	if targetAndSaltMatch[3] != "" {
		targetBracketsDt, err := base64.StdEncoding.DecodeString(targetAndSaltMatch[3])
		if err != nil {
			targetBrackets = targetAndSaltMatch[3]
		} else {
			targetBrackets = string(targetBracketsDt)
		}
	}

	return target, targetBrackets, salt, operation
}

func shortDigest(d digest.Digest) string {
	return d.Hex()[:12]
}
