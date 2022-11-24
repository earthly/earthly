package solvermon

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/util/vertexmeta"
	"github.com/earthly/earthly/util/xcontext"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

// SolverMonitor is a buildkit solver monitor.
type SolverMonitor struct {
	b        *logbus.Bus
	vertices map[digest.Digest]*vertexMonitor
	mu       sync.Mutex
}

// New creates a new SolverMonitor.
func New(b *logbus.Bus) *SolverMonitor {
	return &SolverMonitor{
		b:        b,
		vertices: make(map[digest.Digest]*vertexMonitor),
	}
}

// MonitorProgress processes a channel of buildkit solve statuses.
func (sm *SolverMonitor) MonitorProgress(ctx context.Context, ch chan *client.SolveStatus) error {
	delayedCtx, delayedCancel := context.WithCancel(xcontext.Detach(ctx))
	defer delayedCancel()
	go func() {
		<-ctx.Done()
		// Delay closing to allow any pending messages to be processed.
		// The delay is very high because we expect the buildkit connection
		// to be closed (and hence status channel to be closed) on cancellations
		// anyway. We should be waiting for the full 30 seconds only if there's
		// a bug.
		select {
		case <-delayedCtx.Done():
		case <-time.After(30 * time.Second):
		}
		delayedCancel()
	}()
	for {
		select {
		case <-delayedCtx.Done():
			return errors.Wrap(ctx.Err(), "timed out waiting for status channel to close")
		case status, ok := <-ch:
			if !ok {
				return nil
			}
			err := sm.handleBuildkitStatus(delayedCtx, status)
			if err != nil {
				return err
			}
		}
	}
}

func (sm *SolverMonitor) handleBuildkitStatus(ctx context.Context, status *client.SolveStatus) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	bp := sm.b.Run()
	for _, vertex := range status.Vertexes {
		vm, exists := sm.vertices[vertex.Digest]
		if !exists {
			meta, operation := vertexmeta.ParseFromVertexPrefix(vertex.Name)
			if meta.CanonicalTargetName == "" {
				meta.CanonicalTargetName = meta.TargetName
			}
			var tp *logbus.Target
			if meta.TargetID != "" && meta.TargetName != "" {
				var ok bool
				tp, ok = bp.Target(meta.TargetID)
				if !ok {
					var err error
					tp, err = bp.NewTarget(
						meta.TargetID, meta.TargetName, meta.CanonicalTargetName,
						argsToSlice(meta.OverridingArgs), meta.Platform, meta.Runner)
					if err != nil {
						return err
					}
					tp.SetStart(time.Now())
				}
			}
			push := false // TODO(vladaionescu): Support push.
			cp, err := bp.NewCommand(
				vertex.Digest.String(), operation, meta.TargetID, meta.Platform,
				vertex.Cached, push, meta.Local, meta.SourceLocation,
				meta.RepoGitURL, meta.RepoGitHash, meta.RepoFileRelToRepo)
			if err != nil {
				return err
			}
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
			var status logstream.RunStatus
			switch {
			case vm.isCanceled:
				status = logstream.RunStatus_RUN_STATUS_CANCELED
			case vertex.Error == "" && !vm.isFatalError:
				status = logstream.RunStatus_RUN_STATUS_SUCCESS
			default:
				status = logstream.RunStatus_RUN_STATUS_FAILURE
			}
			vm.cp.SetEnd(*vertex.Completed, status, vm.errorStr)
			if vm.tp != nil {
				// TODO (vladaionescu): The end event is set repeatedly for the
				//                      same target, because we don't know which
				//                      command is the last one for a target.
				//                      This means that some targets can be
				//                      deemed as successful initially, only to be
				//                      overwritten by a failure.
				vm.tp.SetEnd(*vertex.Completed, status, vm.meta.Platform)
			}
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
