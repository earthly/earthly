package builder

import (
	"context"
	"regexp"
	"strings"
	"time"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/logging"
	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

const durationBetweenProgressUpdate = time.Second * 5

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

func (vm *vertexMonitor) printHeader() {
	out := []string{"-->"}
	out = append(out, vm.operation)
	c := vm.console
	if vm.vertex.Cached {
		c = c.WithCached(true)
	}
	c.Printf("%s\n", strings.Join(out, " "))
	vm.headerPrinted = true
}

func (vm *vertexMonitor) shouldPrintProgress(percent int) bool {
	now := time.Now()
	if !vm.headerPrinted {
		return false
	}
	if now.Sub(vm.lastOutput) < durationBetweenProgressUpdate && percent < 100 {
		return false
	}
	if vm.lastPercentage == percent {
		return false
	}
	vm.lastOutput = now
	vm.lastPercentage = percent
	return true
}

func (sm *solverMonitor) monitorProgress(ctx context.Context, ch chan *client.SolveStatus, printDetailed bool) error {
	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				return nil
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
					((printDetailed && !vm.isInternal && (vertex.Cached || vertex.Started != nil)) || vertex.Error != "") {
					vm.printHeader()
					vm.logger.Info("Vertex started or cached")
				}
				if vertex.Error != "" {
					if strings.Contains(vertex.Error, "context canceled: context canceled") {
						if !vm.isInternal {
							vm.console.Printf("WARN: canceled\n")
						}
					} else {
						vm.isError = true
						if strings.Contains(vertex.Error, "executor failed running") {
							vm.console.Warnf("ERROR: Command exited with non-zero code: %s\n", vm.operation)
						} else {
							vm.console.Warnf("ERROR: (%s) %s\n", vm.operation, vertex.Error)
						}
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
				progress := int32(0)
				if vs.Total != 0 {
					progress = int32(100.0 * float32(vs.Current) / float32(vs.Total))
				}
				if vs.Completed != nil {
					progress = 100
				}
				if vm.shouldPrintProgress(int(progress)) {
					logger := vm.logger.
						With("progress", progress).
						With("name", vs.Name)
					if !vm.headerPrinted && printDetailed {
						vm.printHeader()
					}
					logger.Info(vs.ID)
					if printDetailed {
						vm.console.Printf("%s %d%%\n", vs.ID, progress)
					}
				}
			}
			for _, logLine := range ss.Logs {
				vm, ok := sm.vertices[logLine.Vertex]
				if !ok || vm.isInternal {
					// No logging for internal operations.
					continue
				}
				if !vm.headerPrinted && printDetailed {
					vm.printHeader()
				}
				vm.logger.Info(string(logLine.Data))
				if printDetailed {
					vm.console.PrintBytes(logLine.Data)
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
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
