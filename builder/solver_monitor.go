package builder

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/armon/circbuf"
	"github.com/dustin/go-humanize"
	"github.com/earthly/earthly/conslogging"
	"github.com/mattn/go-isatty"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

const (
	durationBetweenSha256ProgressUpdate  = 5 * time.Second
	durationBetweenProgressUpdate        = 3 * time.Second
	durationBetweenProgressUpdateIfSame  = 5 * time.Millisecond
	durationBetweenOpenLineUpdate        = time.Second
	durationBetweenNoOutputUpdates       = 5 * time.Second
	durationBetweenNoOutputUpdatesNoAnsi = 60 * time.Second
	tailErrorBufferSizeBytes             = 80 * 1024 // About as much as 1024 lines of 80 chars each.
)

type vertexMonitor struct {
	vertex         *client.Vertex
	targetStr      string
	targetBrackets string
	meta           map[string]string
	salt           string
	operation      string
	lastProgress   map[string]time.Time
	lastPercentage map[string]int
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

func (vm *vertexMonitor) printHeader() {
	vm.headerPrinted = true
	if vm.operation == "" {
		return
	}
	c := vm.console
	if vm.targetBrackets != "" {
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
	builder = append(builder, fmt.Sprintf("[%s] %s ... %d%%%s\n", progressBar, id, progress, string(ansiEraseRestLine)))
	vm.console.PrintBytes([]byte(strings.Join(builder, "")))
}

func (vm *vertexMonitor) printError() bool {
	if strings.Contains(vm.vertex.Error, "executor failed running") {
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

type solverMonitor struct {
	msgMu                       sync.Mutex
	console                     conslogging.ConsoleLogger
	verbose                     bool
	disableNoOutputUpdates      bool
	vertices                    map[digest.Digest]*vertexMonitor
	saltSeen                    map[string]bool
	lastVertexOutput            *vertexMonitor
	lastOutputWasProgress       bool
	lastOutputWasNoOutputUpdate bool
	timingTable                 map[timingKey]time.Duration
	startTime                   time.Time
	noOutputTicker              *time.Ticker
	noOutputTick                time.Duration
	errVertex                   *vertexMonitor

	mu             sync.Mutex
	success        bool
	ongoing        bool
	printedSuccess bool
}

type timingKey struct {
	targetStr      string
	targetBrackets string
	salt           string
}

func newSolverMonitor(console conslogging.ConsoleLogger, verbose bool, disableNoOutputUpdates bool) *solverMonitor {
	noOutputTick := durationBetweenNoOutputUpdatesNoAnsi
	if ansiSupported {
		noOutputTick = durationBetweenNoOutputUpdates
	}
	return &solverMonitor{
		console:                console,
		verbose:                verbose,
		disableNoOutputUpdates: disableNoOutputUpdates,
		vertices:               make(map[digest.Digest]*vertexMonitor),
		saltSeen:               make(map[string]bool),
		timingTable:            make(map[timingKey]time.Duration),
		startTime:              time.Now(),
		noOutputTicker:         time.NewTicker(noOutputTick),
		noOutputTick:           noOutputTick,
	}
}

func (sm *solverMonitor) monitorProgress(ctx context.Context, ch chan *client.SolveStatus, phaseText string, sideRun bool) (string, error) {
	if !sideRun {
		sm.mu.Lock()
		sm.ongoing = true
		sm.mu.Unlock()
	}
Loop:
	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				break Loop
			}
			err := sm.processStatus(ss)
			if err != nil {
				return "", nil
			}
		case <-sm.noOutputTicker.C:
			err := sm.processNoOutputTick()
			if err != nil {
				return "", nil
			}
		}
	}
	failedVertexOutput := ""
	if !sideRun {
		sm.msgMu.Lock()
		if sm.errVertex != nil {
			if sm.errVertex.tailOutput != nil {
				failedVertexOutput = string(sm.errVertex.tailOutput.Bytes())
			}
			sm.reprintFailure(sm.errVertex, phaseText)
		}
		sm.msgMu.Unlock()
		sm.mu.Lock()
		if sm.success && !sm.printedSuccess {
			sm.lastOutputWasProgress = false
			sm.lastOutputWasNoOutputUpdate = false
			sm.console.PrintSuccess(phaseText)
			sm.printedSuccess = true
		}
		sm.ongoing = false
		sm.mu.Unlock()
		sm.PrintTiming()
		sm.noOutputTicker.Stop()
	}
	return failedVertexOutput, nil
}

func (sm *solverMonitor) processStatus(ss *client.SolveStatus) error {
	sm.msgMu.Lock()
	defer sm.msgMu.Unlock()
	for _, vertex := range ss.Vertexes {
		vm, ok := sm.vertices[vertex.Digest]
		if !ok {
			targetStr, targetBrackets, vertexMetadata, salt, operation := parseVertexName(vertex.Name)
			vm = &vertexMonitor{
				vertex:         vertex,
				targetStr:      targetStr,
				targetBrackets: targetBrackets,
				meta:           vertexMetadata,
				salt:           salt,
				operation:      operation,
				isInternal:     (targetStr == "internal" && !sm.verbose),
				console:        sm.console.WithPrefixAndSalt(targetStr, salt),
				lastPercentage: make(map[string]int),
				lastProgress:   make(map[string]time.Time),
			}
			if vm.meta["@local"] == "true" {
				vm.console = vm.console.WithLocal(true)
			}
			sm.vertices[vertex.Digest] = vm
		}
		vm.vertex = vertex
		if !vm.headerPrinted &&
			((!vm.isInternal && (vertex.Cached || vertex.Started != nil)) || vertex.Error != "") {
			sm.printHeader(vm)
			sm.noOutputTicker.Reset(sm.noOutputTick)
		}
		if vertex.Error != "" {
			if strings.Contains(vertex.Error, "context canceled") {
				if !vm.isInternal {
					vm.console.Printf("WARN: Canceled\n")
					sm.noOutputTicker.Reset(sm.noOutputTick)
				}
			} else {
				vm.isError = vm.printError()
				if sm.errVertex == nil && vm.isError {
					sm.errVertex = vm
				}
				sm.noOutputTicker.Reset(sm.noOutputTick)
			}
		}
		if sm.verbose {
			vm.printTimingInfo()
			sm.recordTiming(vm.targetStr, vm.targetBrackets, vm.salt, vertex)
			sm.noOutputTicker.Reset(sm.noOutputTick)
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
		sm.printProgress(vm, vs.ID, progress)
		sm.noOutputTicker.Reset(sm.noOutputTick)
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
		sm.noOutputTicker.Reset(sm.noOutputTick)
	}
	return nil
}

func (sm *solverMonitor) processNoOutputTick() error {
	if sm.disableNoOutputUpdates {
		return nil
	}
	ongoingBuilder := []string{}
	if sm.lastOutputWasNoOutputUpdate {
		// Overwrite previous line if the previous update was also a no-output update.
		ongoingBuilder = append(ongoingBuilder, string(ansiUp))
	}
	ongoing := []string{}
	now := time.Now()
	for _, vm := range sm.vertices {
		if !vm.isOngoing() {
			continue
		}
		if vm.meta["@interactive"] == "true" {
			// Don't print ongoing updates when an interactive session is ongoing.
			return nil
		}

		col := vm.console.PrefixColor()
		relTime := humanize.RelTime(*vm.vertex.Started, now, "ago", "from now")
		ongoing = append(ongoing, fmt.Sprintf("%s (%s)", col.Sprintf("%s", vm.targetStr), relTime))
	}
	sort.Strings(ongoing) // not entirely correct, but makes the ordering consistent
	var ongoingStr string
	if len(ongoing) > 2 {
		ongoingStr = fmt.Sprintf("%s and %d others", strings.Join(ongoing[:2], ", "), len(ongoing)-2)
	} else {
		ongoingStr = strings.Join(ongoing, ", ")
	}
	ongoingBuilder = append(ongoingBuilder, ongoingStr, string(ansiEraseRestLine))
	sm.console.WithPrefix("ongoing").Printf("%s\n", strings.Join(ongoingBuilder, ""))
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = true
	return nil
}

func (sm *solverMonitor) printOutput(vm *vertexMonitor, data []byte) error {
	sameAsLast := (sm.lastVertexOutput == vm && !sm.lastOutputWasProgress)
	sm.lastVertexOutput = vm
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = false
	return vm.printOutput(data, sameAsLast)
}

func (sm *solverMonitor) printProgress(vm *vertexMonitor, id string, progress int) {
	if vm.shouldPrintProgress(id, progress, sm.verbose, sm.lastOutputWasProgress) {
		if !vm.headerPrinted {
			sm.printHeader(vm)
		}
		vm.printProgress(id, progress, sm.verbose, sm.lastOutputWasProgress)
		sm.lastOutputWasProgress = (progress != 100)
		sm.lastOutputWasNoOutputUpdate = false
	}
}

func (sm *solverMonitor) printHeader(vm *vertexMonitor) {
	seen := sm.saltSeen[vm.salt]
	if !seen {
		sm.saltSeen[vm.salt] = true
	}
	vm.printHeader()
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = false
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

func (sm *solverMonitor) SetSuccess(msg string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.success = true
	if !sm.ongoing {
		sm.lastOutputWasProgress = false
		sm.lastOutputWasNoOutputUpdate = false
		sm.console.PrintSuccess(msg)
		sm.printedSuccess = true
	}
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
	type durationAndKey struct {
		dur time.Duration
		key timingKey
	}
	durs := make([]durationAndKey, 0, len(sm.timingTable))
	for key, dur := range sm.timingTable {
		durs = append(durs, durationAndKey{
			dur: dur,
			key: key,
		})
		total += dur
	}
	sort.Slice(durs, func(i, j int) bool {
		return durs[i].dur > durs[j].dur
	})
	for _, d := range durs {
		sm.console.
			WithPrefixAndSalt(d.key.targetStr, d.key.salt).
			WithMetadataMode(true).
			Printf("(%s) %s\n", d.key.targetBrackets, d.dur)
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

func (sm *solverMonitor) reprintFailure(errVertex *vertexMonitor, phaseText string) {
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = false
	sm.console.Warnf("Repeating the output of the command that caused the failure\n")
	sm.console.PrintFailure(phaseText)
	errVertex.console = errVertex.console.WithFailed(true)
	errVertex.printHeader()
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

var vertexRegexp = regexp.MustCompile("(?s)^\\[([^\\]]*)\\] (.*)$")
var targetAndSaltRegexp = regexp.MustCompile("(?s)^([^\\(]*)(\\(([^\\)]*)\\))? (.*)$")

func parseVertexName(vertexName string) (string, string, map[string]string, string, string) {
	target := ""
	targetBrackets := ""
	operation := ""
	meta := make(map[string]string)
	salt := ""
	if strings.HasPrefix(vertexName, "importing cache manifest") ||
		strings.HasPrefix(vertexName, "exporting cache") {
		return "cache", targetBrackets, meta, "cache", vertexName
	}
	match := vertexRegexp.FindStringSubmatch(vertexName)
	if len(match) < 2 {
		return "internal", targetBrackets, meta, "internal", vertexName
	}
	targetAndSalt := match[1]
	operation = match[2]
	targetAndSaltMatch := targetAndSaltRegexp.FindStringSubmatch(targetAndSalt)
	if targetAndSaltMatch == nil {
		return targetAndSalt, targetBrackets, meta, targetAndSalt, operation
	}
	target = targetAndSaltMatch[1]
	salt = targetAndSaltMatch[len(targetAndSaltMatch)-1]
	if salt == "" {
		salt = targetAndSalt
	}
	if targetAndSaltMatch[3] != "" {
		targetBracketsParts := strings.Split(targetAndSaltMatch[3], " ")
		targetBracketsBuilder := make([]string, 0, len(targetBracketsParts))
		for _, part := range targetBracketsParts {
			kvParts := strings.SplitN(part, "=", 2)
			if len(kvParts) != 2 {
				// Handle an error better here?
				continue
			}
			key := kvParts[0]
			var value string
			valueDt, err := base64.StdEncoding.DecodeString(kvParts[1])
			if err != nil {
				// Handle an error better here?
				value = kvParts[1]
			} else {
				value = string(valueDt)
			}
			if strings.HasPrefix(key, "@") {
				meta[key] = value
			} else {
				targetBracketsBuilder = append(targetBracketsBuilder, fmt.Sprintf("%s=%s", key, value))
			}
		}
		platform := meta["@platform"]
		if platform != "" {
			targetBracketsBuilder = append(
				[]string{fmt.Sprintf("platform=%s", platform)}, targetBracketsBuilder...)
		}
		targetBrackets = strings.Join(targetBracketsBuilder, " ")
	}
	return target, targetBrackets, meta, salt, operation
}

func shortDigest(d digest.Digest) string {
	return d.Hex()[:12]
}

var progressChars = []string{
	" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█",
}

func progressBar(progress, width int) string {
	if progress > 100 {
		progress = 100
	}
	if progress < 0 {
		progress = 0
	}
	builder := make([]string, 0, width)
	fullChars := progress * width / 100
	blankChars := width - fullChars - 1
	deltaProgress := ((progress * width) % 100) * len(progressChars) / 100
	for i := 0; i < fullChars; i++ {
		builder = append(builder, progressChars[len(progressChars)-1])
	}
	if progress != 100 {
		builder = append(builder, progressChars[deltaProgress])
	}
	for i := 0; i < blankChars; i++ {
		builder = append(builder, progressChars[0])
	}
	return strings.Join(builder, "")
}
