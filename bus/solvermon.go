package bus

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/earthly/earthly/outmon"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
)

const tailErrorBufferSizeBytes = 80 * 1024 // About as much as 1024 lines of 80 chars each.

// SolverMonitor is a buildkit solver monitor.
type SolverMonitor struct {
	b        *Bus
	vertices map[digest.Digest]*vertexMonitor
	mu       sync.Mutex
}

func newSolverMonitor(b *Bus) *SolverMonitor {
	return &SolverMonitor{
		b:        b,
		vertices: make(map[digest.Digest]*vertexMonitor),
	}
}

// MonitorProgress processes a channel of buildkit solve statuses.
func (sm *SolverMonitor) MonitorProgress(ctx context.Context, ch chan *client.SolveStatus) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case status := <-ch:
			return sm.handleBuildkitStatus(ctx, status)
		}
	}
}

func (sm *SolverMonitor) handleBuildkitStatus(ctx context.Context, status *client.SolveStatus) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	bp := sm.b.Printer()
	for _, vertex := range status.Vertexes {
		vm, exists := sm.vertices[vertex.Digest]
		if !exists {
			meta, operation := outmon.ParseFromVertexPrefix(vertex.Name)
			tp := bp.TargetPrinter(
				meta.TargetID, meta.TargetName, meta.CanonicalTargetName, argsToSlice(meta.OverridingArgs), meta.Platform)
			_, cp := tp.NextCommandPrinter(operation, vertex.Cached, false, meta.Local, meta.SourceLocation)
			vm = &vertexMonitor{
				vertex:    vertex,
				meta:      meta,
				operation: operation,
				cp:        cp,
			}
			sm.vertices[vertex.Digest] = vm
		}
		vm.vertex = vertex
		if vertex.Started != nil {
			vm.cp.SetStart(*vertex.Started)
		}
		if vertex.Error != "" {
			vm.parseError()
		}
		if vertex.Completed != nil {
			success := (vertex.Error == "" && !vm.isFatalError && !vm.isCanceled)
			vm.cp.SetEnd(*vertex.Completed, success, vm.isCanceled, vm.errorStr)
			if vm.isFatalError {
				// Run this at the end so that we capture any additional log lines.
				defer func(end time.Time, targetID string, index int32, errorStr string) {
					output := []byte{}
					if vm.tailOutput != nil {
						output = vm.tailOutput.Bytes()
					}
					bp.SetFatalError(end, targetID, true, index, output, errorStr)
				}(*vertex.Completed, vm.meta.TargetID, vm.cp.Index(), vm.errorStr)
			}
		}
	}
	for _, vs := range status.Statuses {
		vm, exists := sm.vertices[vs.Vertex]
		if !exists {
			continue
		}
		progress := int32(0)
		if vs.Total != 0 {
			progress = int32(100.0 * float32(vs.Current) / float32(vs.Total))
		}
		if vs.Completed != nil {
			progress = 100
		}
		vm.cp.SetProgress(progress)
	}
	for _, logLine := range status.Logs {
		vm, exists := sm.vertices[logLine.Vertex]
		if !exists {
			continue
		}
		_, err := vm.Write(logLine.Data, logLine.Timestamp, logLine.Stream)
		if err != nil {
			return err
		}
	}
	return nil
}

func argsToSlice(args map[string]string) []string {
	var argsSlice []string
	for k, v := range args {
		argsSlice = append(argsSlice, k+"="+v)
	}
	sort.StringSlice(argsSlice).Sort()
	return argsSlice
}
