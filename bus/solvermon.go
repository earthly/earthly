package bus

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/earthly/earthly/util/vertexmeta"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
)

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
	returnedCh := make(chan struct{})
	defer close(returnedCh)
	closedCh := make(chan struct{})
	go func() {
		<-ctx.Done()
		// Delay closing to allow any pending messages
		// to be processed.
		select {
		case <-returnedCh:
		case <-time.After(5 * time.Second):
		}
		close(closedCh)
	}()
	for {
		select {
		case <-closedCh:
			return ctx.Err()
		case status, ok := <-ch:
			if !ok {
				return nil
			}
			err := sm.handleBuildkitStatus(ctx, status)
			if err != nil {
				return err
			}
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
			meta, operation := vertexmeta.ParseFromVertexPrefix(vertex.Name)
			if meta.CanonicalTargetName == "" {
				meta.CanonicalTargetName = meta.TargetName
			}
			var tp *TargetPrinter
			if meta.TargetID != "" && meta.TargetName != "" {
				var ok bool
				tp, ok = bp.TargetPrinter(meta.TargetID)
				if !ok {
					tp = bp.NewTargetPrinter(
						meta.TargetID, meta.TargetName, meta.CanonicalTargetName,
						argsToSlice(meta.OverridingArgs), meta.Platform)
					// TODO(vladaionescu): All the target printers should get
					//                     SetStart and SetEnd appropriately.
				}
			}
			push := false // TODO(vladaionescu): Support push.
			cp := bp.NewCommandPrinter(
				vertex.Digest.String(), operation, meta.TargetID, meta.Platform,
				vertex.Cached, push, meta.Local, meta.SourceLocation,
				meta.RepoGitURL, meta.RepoGitHash, meta.RepoFileRelToRepo)
			vm = &vertexMonitor{
				vertex:    vertex,
				meta:      meta,
				operation: operation,
				tp:        tp,
				cp:        cp,
			}
			sm.vertices[vertex.Digest] = vm
		}
		vm.vertex = vertex
		if vertex.Cached {
			vm.cp.SetCached(true)
		}
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
				defer bp.SetFatalError(
					*vertex.Completed, vm.meta.TargetID, vm.vertex.Digest.String(),
					vm.fatalErrorType, vm.errorStr)
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
