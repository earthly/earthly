package outmon

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/logbus/formatter"
	"github.com/earthly/earthly/util/buildkitutil"
	"github.com/earthly/earthly/util/vertexmeta"
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

// SolverMonitor is an object that monitors for status updates from a buildkit solve
// and prints them to the console.
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

// MonitorProgress consumes progress messages from a solve statue channel and prints them to the console.
func (sm *SolverMonitor) MonitorProgress(ctx context.Context, ch chan *client.SolveStatus, phaseText string, sideRun bool, bkClient *client.Client) (string, error) {
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
			err := sm.processNoOutputTick(ctx, bkClient)
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
			meta, operation := vertexmeta.ParseFromVertexPrefix(vertex.Name)
			vm = &vertexMonitor{
				vertex:         vertex,
				meta:           meta,
				operation:      operation,
				console:        sm.console.WithPrefixAndSalt(meta.TargetName, meta.Salt()),
				lastPercentage: make(map[string]int),
				lastProgress:   make(map[string]time.Time),
			}
			if vm.meta.Local {
				vm.console = vm.console.WithLocal(true)
			}
			sm.vertices[vertex.Digest] = vm
		}
		vm.vertex = vertex
		if !vm.headerPrinted &&
			((!vm.meta.Internal && (vertex.Cached || vertex.Started != nil)) || vertex.Error != "") {
			sm.printHeader(vm)
			sm.noOutputTicker.Reset(sm.noOutputTick)
		}
		if vertex.Error != "" {
			fmt.Printf("%v got a vertex error %v\n", time.Now(), vertex.Error)
			if strings.Contains(vertex.Error, "context canceled") {
				if !vm.meta.Internal {
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
			sm.recordTiming(vm, vertex)
			sm.noOutputTicker.Reset(sm.noOutputTick)
		}

		vm.reportStatusToConsole()
		vm.reportResultToConsole()
	}
	for _, vs := range ss.Statuses {
		vm, ok := sm.vertices[vs.Vertex]
		if !ok || vm.meta.Internal {
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
		if !ok || vm.meta.Internal {
			// No logging for internal operations.
			continue
		}
		if !vm.headerPrinted {
			sm.printHeader(vm)
		}
		if logLine.Stream == formatter.BuildkitStatsStream {
			// --exec-stats requires logbus to be enabled
			continue
		}
		err := sm.printOutput(vm, logLine.Data)
		if err != nil {
			return err
		}
		sm.noOutputTicker.Reset(sm.noOutputTick)
	}
	return nil
}

func (sm *SolverMonitor) processNoOutputTick(ctx context.Context, bkClient *client.Client) error {
	ongoingCons := sm.console.WithPrefix("ongoing")
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
		if vm.meta.Interactive {
			// Don't print ongoing updates when an interactive session is ongoing.
			return nil
		}

		col := vm.console.PrefixColor()
		relTime := humanize.RelTime(*vm.vertex.Started, now, "ago", "from now")
		ongoing = append(ongoing, fmt.Sprintf("%s (%s)", col.Sprintf("%s", vm.meta.TargetName), relTime))
	}
	var ongoingStr string
	warn := false
	defer func() {
		// Note: This part assumes that we are still under lock.
		ongoingBuilder = append(ongoingBuilder, ongoingStr, string(ansiEraseRestLine))
		outStr := strings.Join(ongoingBuilder, "")
		if warn {
			ongoingCons.Warnf("%s\n", outStr)
		} else {
			ongoingCons.Printf("%s\n", outStr)
		}
		sm.lastOutputWasProgress = false
		sm.lastOutputWasNoOutputUpdate = true
	}()

	if len(ongoing) != 0 {
		sort.Strings(ongoing) // not entirely correct, but makes the ordering consistent
		if len(ongoing) > 2 {
			ongoingStr = fmt.Sprintf("%s and %d others", strings.Join(ongoing[:2], ", "), len(ongoing)-2)
		} else {
			ongoingStr = strings.Join(ongoing, ", ")
		}
		return nil
	}

	// Nothing running, but also no output taking place. We're just sitting,
	// waiting for buildkit to make progress. Let's check if buildkit is
	// overwhelmed and report accordingly.
	workers, err := bkClient.ListWorkers(ctx)
	if err != nil {
		ongoingStr = fmt.Sprintf("error getting buildkit worker info: %v", err)
		warn = true
		return nil // no need to crash
	}
	numOtherSessions := -1 // default to unknown (since Info is a newer call and might not be implemented)
	if info, err := bkClient.Info(ctx); err == nil {
		numOtherSessions = info.NumSessions - 1 // don't count current session
	}

	if len(workers) == 0 {
		ongoingStr = "error getting buildkit worker info: no workers"
		warn = true
		return nil // no need to crash
	}
	workerInfo := workers[0]
	load := workerInfo.ParallelismCurrent + workerInfo.ParallelismWaiting
	switch {
	case workerInfo.ParallelismWaiting > 5:
		ongoingStr = fmt.Sprintf(
			"Waiting... Buildkit is currently under heavy load (%s)", buildkitutil.FormatUtilization(numOtherSessions, load, workerInfo.ParallelismMax))
		warn = true
	case workerInfo.ParallelismWaiting > 0:
		ongoingStr = fmt.Sprintf(
			"Waiting... Buildkit is currently under significant load (%s)", buildkitutil.FormatUtilization(numOtherSessions, load, workerInfo.ParallelismMax))
	default:
		ongoingStr = fmt.Sprintf(
			"Waiting on Buildkit... (%s)", buildkitutil.FormatUtilization(numOtherSessions, load, workerInfo.ParallelismMax))
	}
	if workerInfo.GCAnalytics.CurrentStartTime != nil {
		d := now.Sub(*workerInfo.GCAnalytics.CurrentStartTime).Round(time.Second)
		if d <= 5*time.Minute {
			ongoingStr += fmt.Sprintf(" GC (%v ago)", d)
		} else {
			ongoingStr += fmt.Sprintf(" GC is slow (%v ago)", d)
			warn = true
		}
	}
	if workerInfo.GCAnalytics.AllTimeMaxDuration > 5*time.Minute {
		ongoingStr += fmt.Sprintf(
			" GCs historically slow (max %v)",
			workerInfo.GCAnalytics.AllTimeMaxDuration.Round(time.Second))
		warn = true
	}
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
	seen := sm.saltSeen[vm.meta.Salt()]
	if !seen {
		sm.saltSeen[vm.meta.Salt()] = true
	}
	vm.printHeader(sm.verbose)
	sm.lastOutputWasProgress = false
	sm.lastOutputWasNoOutputUpdate = false
}

func (sm *SolverMonitor) recordTiming(vm *vertexMonitor, vertex *client.Vertex) {
	if vertex.Started == nil || vertex.Completed == nil {
		return
	}
	dur := vertex.Completed.Sub(*vertex.Started)
	if dur == 0 {
		return
	}
	key := timingKey{
		targetStr:      vm.meta.TargetName,
		targetBrackets: vm.meta.OverridingArgsString(),
		salt:           vm.meta.Salt(),
	}
	sm.timingTable[key] += dur
}

// PrintTiming prints the accumulated timing information.
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
	errVertex.printHeader(sm.verbose)
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
