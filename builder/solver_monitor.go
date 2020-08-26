package builder

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/logging"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

const (
	durationBetweenProgressUpdate = time.Second * 5
	tailErrorBufferSizeBytes      = 80 * 1024 // About as much as 1024 lines of 80 chars each.
)

type vertexMonitor struct {
	vertex         *client.Vertex
	targetStr      string
	salt           string
	operation      string
	lastOutput     time.Time
	lastPercentage int
	logger         logging.Logger
	console        conslogging.ConsoleLogger
	headerPrinted  bool
	isInternal     bool
	isError        bool
	tailOutput     *circbuf.Buffer
}

func (vm *vertexMonitor) printHeader() {
	vm.headerPrinted = true
	if vm.operation == "" {
		return
	}
	out := []string{"-->"}
	out = append(out, vm.operation)
	c := vm.console
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

func (vm *vertexMonitor) printOutput(output []byte) error {
	vm.console.PrintBytes(output)
	if vm.tailOutput == nil {
		var err error
		vm.tailOutput, err = circbuf.NewBuffer(tailErrorBufferSizeBytes)
		if err != nil {
			return errors.Wrap(err, "allocate buffer for output")
		}
	}
	_, err := vm.tailOutput.Write(output)
	if err != nil {
		return errors.Wrap(err, "write to in-memory output buffer")
	}
	return nil
}

func (vm *vertexMonitor) printError() {
	if strings.Contains(vm.vertex.Error, "executor failed running") {
		vm.console.Warnf("ERROR: Command exited with non-zero code: %s\n", vm.operation)
	} else {
		vm.console.Warnf("ERROR: (%s) %s\n", vm.operation, vm.vertex.Error)
	}
}

type solverMonitor struct {
	console conslogging.ConsoleLogger

	vertices map[digest.Digest]*vertexMonitor
}

func newSolverMonitor(console conslogging.ConsoleLogger) *solverMonitor {
	return &solverMonitor{
		console:  console,
		vertices: make(map[digest.Digest]*vertexMonitor),
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
					targetStr, salt, operation := parseVertexName(vertex.Name)
					vertexLogger := logging.GetLogger(ctx).
						With("target", targetStr).
						With("vertex", shortDigest(vertex.Digest)).
						With("cached", vertex.Cached).
						With("operation", operation)
					vm = &vertexMonitor{
						vertex:     vertex,
						targetStr:  targetStr,
						salt:       salt,
						operation:  operation,
						logger:     vertexLogger,
						isInternal: (targetStr == "internal"),
						console:    sm.console.WithPrefixAndSalt(targetStr, salt),
					}
					sm.vertices[vertex.Digest] = vm
				}
				vm.vertex = vertex
				if !vm.headerPrinted &&
					((!vm.isInternal && (vertex.Cached || vertex.Started != nil)) || vertex.Error != "") {
					vm.printHeader()
					vm.logger.Info("Vertex started or cached")
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
					vm.logger.Error(errors.New(vertex.Error))
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
					logger := vm.logger.
						With("progress", progress).
						With("name", vs.Name)
					if !vm.headerPrinted {
						vm.printHeader()
					}
					logger.Info(vs.ID)
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
					vm.printHeader()
				}
				vm.logger.Info(string(logLine.Data))
				err := vm.printOutput(logLine.Data)
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

func (sm *solverMonitor) reprintFailure(errVertex *vertexMonitor) {
	sm.console.Warnf("Repeating the output of the command that caused the failure\n")
	sm.console.PrintFailure()
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

var bracketsRegexp = regexp.MustCompile("^\\[([^\\]]*)\\] (.*)$")

func parseVertexName(vertexName string) (string, string, string) {
	target := ""
	operation := ""
	salt := ""
	match := bracketsRegexp.FindStringSubmatch(vertexName)
	if len(match) < 2 {
		return target, salt, operation
	}
	targetAndSalt := match[1]
	targetAndSaltSlice := strings.SplitN(targetAndSalt, " ", 2)
	if len(targetAndSaltSlice) == 2 {
		target = targetAndSaltSlice[0]
		salt = targetAndSaltSlice[1]
	} else {
		target = targetAndSalt
	}
	if len(match) < 3 {
		return target, salt, operation
	}
	operation = match[2]
	return target, salt, operation
}

func shortDigest(d digest.Digest) string {
	return d.Hex()[:12]
}
