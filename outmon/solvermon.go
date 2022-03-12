package outmon

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/earthly/earthly/conslogging"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
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

type SolverMonitor struct {
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

	mu      sync.Mutex
	ongoing bool
}

type timingKey struct {
	targetStr      string
	targetBrackets string
	salt           string
}

// NewSolverMonitor retuns a new solver monitor.
func NewSolverMonitor(console conslogging.ConsoleLogger, verbose bool, disableNoOutputUpdates bool) *SolverMonitor {
	noOutputTick := durationBetweenNoOutputUpdatesNoAnsi
	if ansiSupported {
		noOutputTick = durationBetweenNoOutputUpdates
	}
	return &SolverMonitor{
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

// Monitor progress consumes progress messages from a solve statue channel and prints them to the console.
func (sm *SolverMonitor) MonitorProgress(ctx context.Context, ch chan *client.SolveStatus, phaseText string, sideRun bool) (string, error) {
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
				return "", err
			}
		case <-sm.noOutputTicker.C:
			err := sm.processNoOutputTick()
			if err != nil {
				return "", err
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
		sm.ongoing = false
		sm.mu.Unlock()
		sm.PrintTiming()
		sm.noOutputTicker.Stop()
	}
	return failedVertexOutput, nil
}

func (sm *SolverMonitor) processStatus(ss *client.SolveStatus) error {
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
					vm.isCanceled = true
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

		vm.reportStatusToConsole()
		vm.reportResultToConsole()
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

func (sm *SolverMonitor) processNoOutputTick() error {
	sm.msgMu.Lock()
	defer sm.msgMu.Unlock()
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

func (sm *SolverMonitor) printOutput(vm *vertexMonitor, data []byte) error {
	sameAsLast := (sm.lastVertexOutput == vm && !sm.lastOutputWasProgress)
	sm.lastVertexOutput = vm
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = false
	return vm.printOutput(data, sameAsLast)
}

func (sm *SolverMonitor) printProgress(vm *vertexMonitor, id string, progress int) {
	if vm.shouldPrintProgress(id, progress, sm.verbose, sm.lastOutputWasProgress) {
		if !vm.headerPrinted {
			sm.printHeader(vm)
		}
		vm.printProgress(id, progress, sm.verbose, sm.lastOutputWasProgress)
		sm.lastOutputWasProgress = (progress != 100)
		sm.lastOutputWasNoOutputUpdate = false
	}
}

func (sm *SolverMonitor) printHeader(vm *vertexMonitor) {
	seen := sm.saltSeen[vm.salt]
	if !seen {
		sm.saltSeen[vm.salt] = true
	}
	vm.printHeader()
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = false
}

func (sm *SolverMonitor) recordTiming(targetStr, targetBrackets, salt string, vertex *client.Vertex) {
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

func (sm *SolverMonitor) PrintTiming() {
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
		Printf("Total (real)\t%s\n", time.Since(sm.startTime))
}

func (sm *SolverMonitor) reprintFailure(errVertex *vertexMonitor, phaseText string) {
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = false
	sm.console.PrintFailure(phaseText)
	sm.console.Warnf("Repeating the output of the command that caused the failure\n")
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

var vertexRegexp = regexp.MustCompile(`(?s)^\[([^\]]*)\] (.*)$`)
var targetAndSaltRegexp = regexp.MustCompile(`(?s)^([^\(]*)(\(([^\)]*)\))? (.*)$`)

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
