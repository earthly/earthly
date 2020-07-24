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

var durationBetweenProgressUpdate = time.Second * 5

type lastProgress struct {
	lastOutput     time.Time
	lastPercentage int
}

type solverMonitor struct {
	console conslogging.ConsoleLogger

	lastProgress map[digest.Digest]lastProgress
}

func newSolverMonitor(console conslogging.ConsoleLogger) *solverMonitor {
	return &solverMonitor{
		console:      console,
		lastProgress: map[digest.Digest]lastProgress{},
	}
}

func (s *solverMonitor) shouldPrint(vert digest.Digest, percent int) bool {
	now := time.Now()
	lp, ok := s.lastProgress[vert]
	if !ok {
		s.lastProgress[vert] = lastProgress{
			lastOutput:     now,
			lastPercentage: percent,
		}
		return true
	}
	if now.Sub(lp.lastOutput) < durationBetweenProgressUpdate && percent < 100 {
		return false
	}
	if lp.lastPercentage == percent {
		return false
	}

	s.lastProgress[vert] = lastProgress{
		lastOutput:     now,
		lastPercentage: percent,
	}
	return true
}

func (s *solverMonitor) monitorProgress(ctx context.Context, ch chan *client.SolveStatus, printDetailed bool) error {
	vertexLoggers := make(map[digest.Digest]logging.Logger)
	vertexConsoles := make(map[digest.Digest]conslogging.ConsoleLogger)
	vertices := make(map[digest.Digest]*client.Vertex)
	introducedVertex := make(map[digest.Digest]bool)
	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				return nil
			}
			for _, vertex := range ss.Vertexes {
				if strings.HasPrefix(vertex.Name, "[internal]") {
					// No logging for internal operations.
					continue
				}
				targetStr, operation := parseVertexName(vertex.Name)
				logger := logging.GetLogger(ctx).
					With("target", targetStr).
					With("vertex", shortDigest(vertex.Digest)).
					With("cached", vertex.Cached).
					With("operation", operation)
				vertexLoggers[vertex.Digest] = logger
				targetConsole := s.console.WithPrefix(targetStr)
				vertexConsoles[vertex.Digest] = targetConsole
				vertices[vertex.Digest] = vertex
				if !introducedVertex[vertex.Digest] && (vertex.Cached || vertex.Started != nil) {
					introducedVertex[vertex.Digest] = true
					if printDetailed {
						printVertex(vertex, targetConsole)
					}
					logger.Info("Vertex started or cached")
				}
				if vertex.Error != "" {
					if !introducedVertex[vertex.Digest] {
						introducedVertex[vertex.Digest] = true
						if printDetailed {
							printVertex(vertex, targetConsole)
						}
					}
					if strings.Contains(vertex.Error, "context canceled: context canceled") {
						targetConsole.Printf("WARN: (%s) canceled\n", operation)
					} else {
						if printDetailed {
							targetConsole.Warnf("ERROR: (%s) %s\n", operation, vertex.Error)
						}
					}
					logger.Error(errors.New(vertex.Error))
				}
			}
			for _, vs := range ss.Statuses {
				vertex, found := vertices[vs.Vertex]
				if !found {
					// No logging for internal operations.
					continue
				}
				logger := vertexLoggers[vs.Vertex]
				targetConsole := vertexConsoles[vs.Vertex]
				progress := int32(0)
				if vs.Total != 0 {
					progress = int32(100.0 * float32(vs.Current) / float32(vs.Total))
				}
				if vs.Completed != nil {
					progress = 100
				}
				if s.shouldPrint(vs.Vertex, int(progress)) {
					logger = logger.
						With("progress", progress).
						With("name", vs.Name)
					if !introducedVertex[vertex.Digest] {
						introducedVertex[vertex.Digest] = true
						if printDetailed {
							printVertex(vertex, targetConsole)
						}
					}
					logger.Info(vs.ID)
					if printDetailed {
						targetConsole.Printf("%s %d%%\n", vs.ID, progress)
					}
				}
			}
			for _, logLine := range ss.Logs {
				vertex, found := vertices[logLine.Vertex]
				if !found {
					// No logging for internal operations.
					continue
				}
				logger := vertexLoggers[logLine.Vertex]
				targetConsole := vertexConsoles[logLine.Vertex]
				if !introducedVertex[logLine.Vertex] {
					introducedVertex[logLine.Vertex] = true
					printVertex(vertex, targetConsole)
				}
				logger.Info(string(logLine.Data))
				targetConsole.PrintBytes(logLine.Data)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func printVertex(vertex *client.Vertex, console conslogging.ConsoleLogger) {
	_, operation := parseVertexName(vertex.Name)
	out := []string{"-->"}
	out = append(out, operation)
	c := console
	if vertex.Cached {
		c = c.WithCached(true)
	}
	c.Printf("%s\n", strings.Join(out, " "))
}

var bracketsRegexp = regexp.MustCompile("^\\[([^\\]]*)\\] (.*)$")

func parseVertexName(vertexName string) (string, string) {
	target := ""
	operation := ""
	match := bracketsRegexp.FindStringSubmatch(vertexName)
	if len(match) < 2 {
		return target, operation
	}
	target = match[1]
	if len(match) < 3 {
		return target, operation
	}
	operation = match[2]
	return target, operation
}

func shortDigest(d digest.Digest) string {
	return d.Hex()[:12]
}
