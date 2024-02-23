package solvermon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/earthly/cloud-api/logstream"
	"github.com/earthly/earthly/logbus"
	"github.com/earthly/earthly/util/statsstreamparser"
	"github.com/earthly/earthly/util/stringutil"
	"github.com/earthly/earthly/util/vertexmeta"
	"github.com/earthly/earthly/util/xcontext"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

// SolverMonitor is a buildkit solver monitor.
type SolverMonitor struct {
	b        *logbus.Bus
	digests  map[digest.Digest]string  // digest -> cmdID
	vertices map[string]*vertexMonitor // cmdID -> vertexMonitor
	mu       sync.Mutex
}

// New creates a new SolverMonitor.
func New(b *logbus.Bus) *SolverMonitor {
	return &SolverMonitor{
		b:        b,
		digests:  make(map[digest.Digest]string),
		vertices: make(map[string]*vertexMonitor),
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
		meta, operation := vertexmeta.ParseFromVertexPrefix(vertex.Name)
		var cmdID string
		createCmd := true
		switch {
		case meta.TargetName == "context":
			cmdID = operation
		case meta.CommandID != "":
			// If the command ID is set, the Logbus command is guaranteed to
			// have been created by Earthly in the converter ahead of time.
			cmdID = meta.CommandID
			createCmd = false
		default:
			cmdID = vertex.Digest.String()
		}
		vm, exists := sm.vertices[cmdID]
		if exists {
			sm.digests[vertex.Digest] = cmdID
		} else {
			category := meta.TargetName
			if meta.Internal {
				category = fmt.Sprintf("internal %s", category)
			}
			var cp *logbus.Command
			// Operations initiated from Earthly have created Logbus commands
			// ahead-of-time. Others may originate from BuildKit, so we'll have
			// to create a command at this point.
			if createCmd {
				var err error
				cp, err = bp.NewCommand(
					cmdID, operation, meta.TargetID, category, meta.Platform,
					vertex.Cached, meta.Local, meta.Interactive, meta.SourceLocation,
					meta.RepoGitURL, meta.RepoGitHash, meta.RepoFileRelToRepo)
				if err != nil {
					return err
				}
			} else {
				var ok bool
				cp, ok = bp.Command(cmdID)
				if !ok {
					// Note: if we receive a vertex with a full command ID that
					// does not exist in this process, it may have originated
					// from another Earthly process. It should be safe to
					// ignore, in this case.
					continue
				}
				cp.SetName(operation) // Command created prior may not have a full name.
			}
			vm = &vertexMonitor{
				vertex:    vertex,
				meta:      meta,
				operation: operation,
				cp:        cp,
				ssp:       statsstreamparser.New(),
			}
			sm.vertices[cmdID] = vm
			sm.digests[vertex.Digest] = cmdID
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
			if vm.isFatalError {
				// Run this at the end so that we capture any additional log lines.
				defer bp.SetFatalError(
					*vertex.Completed, vm.meta.TargetID, cmdID,
					vm.fatalErrorType, stringutil.ScrubCredentialsAll(vm.errorStr))
			}
		}
	}
	for _, vs := range status.Statuses {
		cmdID, exists := sm.digests[vs.Vertex]
		if !exists {
			continue
		}
		vm := sm.vertices[cmdID]
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
		cmdID, exists := sm.digests[logLine.Vertex]
		if !exists {
			continue
		}
		vm := sm.vertices[cmdID]
		logLine.Data = []byte(stringutil.ScrubCredentialsAll((string(logLine.Data))))
		_, err := vm.Write(logLine.Data, logLine.Timestamp, logLine.Stream)
		if err != nil {
			return err
		}
	}
	return nil
}
