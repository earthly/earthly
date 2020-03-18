package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/solver/pb"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
	"github.com/vladaionescu/earthly/conslogging"
	"github.com/vladaionescu/earthly/earthfile2llb/image"
	"github.com/vladaionescu/earthly/llbutil"
	"github.com/vladaionescu/earthly/logging"
	"golang.org/x/sync/errgroup"
)

type solver struct {
	bkClient    *client.Client
	console     conslogging.ConsoleLogger
	attachables []session.Attachable
	enttlmnts   []entitlements.Entitlement
}

func (s *solver) solveDocker(ctx context.Context, localDirs map[string]string, state llb.State, img *image.Image, dockerTag string, push bool) error {
	dt, err := state.Marshal(llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	// TODO: Maybe add some buffering?
	pipeR, pipeW := io.Pipe()
	solveOpt, err := s.newSolveOptDocker(img, dockerTag, localDirs, pipeW)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Solve(ctx, dt, *solveOpt, ch)
		if err != nil {
			return errors.Wrap(err, "solve")
		}
		logging.GetLogger(ctx).Info("Solve successful")
		return nil
	})
	eg.Go(func() error {
		return s.monitorProgressBasic(ctx, ch)
	})
	eg.Go(func() error {
		defer pipeR.Close()
		err := loadDockerTar(ctx, pipeR)
		if err != nil {
			return errors.Wrapf(err, "load docker tar for %s", dockerTag)
		}
		logging.GetLogger(ctx).Info("Docker load success")
		if push {
			err := pushDockerImage(ctx, dockerTag)
			if err != nil {
				return err
			}
			logging.GetLogger(ctx).Info("Docker push success")
		}
		return nil
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Close read pipe on cancels, otherwise the whole thing hangs.
				pipeR.Close()
			}
		}
	}()
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "build error group")
	}
	return nil
}

func (s *solver) solveDockerTar(ctx context.Context, localDirs map[string]string, state llb.State, img *image.Image, dockerTag string, outFile string) error {
	dt, err := state.Marshal(llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	pipeR, pipeW := io.Pipe()
	solveOpt, err := s.newSolveOptDocker(img, dockerTag, localDirs, pipeW)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Solve(ctx, dt, *solveOpt, ch)
		if err != nil {
			return errors.Wrap(err, "solve")
		}
		logging.GetLogger(ctx).Info("Solve successful")
		return nil
	})
	eg.Go(func() error {
		return s.monitorProgressBasic(ctx, ch)
	})
	eg.Go(func() error {
		file, err := os.Create(outFile)
		if err != nil {
			return errors.Wrapf(err, "open file %s for writing", outFile)
		}
		defer file.Close()
		bufFile := bufio.NewWriter(file)
		defer bufFile.Flush()
		buf := make([]byte, 1024)
		for {
			n, err := pipeR.Read(buf)
			if err != nil && err != io.EOF {
				return errors.Wrap(err, "pipe read")
			}
			if err == io.EOF {
				break
			}
			_, err = bufFile.Write(buf[:n])
			if err != nil {
				return errors.Wrap(err, "write chunk to file")
			}
		}
		return nil
	})
	go func() {
		for {
			select {
			case <-ctx.Done():
				// Close read pipe on cancels, otherwise the whole thing hangs.
				pipeR.Close()
			}
		}
	}()
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "build error group")
	}
	return nil
}

func (s *solver) solveArtifacts(ctx context.Context, localDirs map[string]string, state llb.State, outDir string) error {
	dt, err := state.Marshal(llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	solveOpt, err := s.newSolveOptArtifacts(outDir, localDirs)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Solve(ctx, dt, *solveOpt, ch)
		if err != nil {
			return errors.Wrap(err, "solve")
		}
		logging.GetLogger(ctx).Info("Solve successful")
		return nil
	})
	eg.Go(func() error {
		return s.monitorProgressBasic(ctx, ch)
	})
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "build error group")
	}
	return nil
}

func (s *solver) solveSideEffects(ctx context.Context, localDirs map[string]string, state llb.State, printDetailed bool) error {
	dt, err := state.Marshal(llb.Platform(llbutil.TargetPlatform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	dtBytes, err := json.Marshal(dt)
	if err != nil {
		return errors.Wrap(err, "json marshal of state")
	}
	var ops []*pb.Op
	for _, opDef := range dt.Def {
		var op pb.Op
		err = proto.Unmarshal(opDef, &op)
		if err != nil {
			return errors.Wrap(err, "proto unmarshal of op")
		}
		ops = append(ops, &op)
	}
	opsBytes, err := json.Marshal(&ops)
	if err != nil {
		return errors.Wrap(err, "json marshal of ops")
	}
	logging.GetLogger(ctx).
		With("ops", string(opsBytes)).
		With("dt", string(dtBytes)).
		Debug("Side effectsLLB")
	solveOpt, err := s.newSolveOptSideEffects(localDirs)
	if err != nil {
		return errors.Wrap(err, "new solve opt")
	}
	ch := make(chan *client.SolveStatus)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		var err error
		_, err = s.bkClient.Solve(ctx, dt, *solveOpt, ch)
		if err != nil {
			return errors.Wrap(err, "solve")
		}
		logging.GetLogger(ctx).Info("Solve successful")
		return nil
	})
	eg.Go(func() error {
		if printDetailed {
			return s.monitorProgressDetailed(ctx, ch)
		}
		return s.monitorProgressBasic(ctx, ch)
	})
	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "build error group")
	}
	return nil
}

func (s *solver) monitorProgressDetailed(ctx context.Context, ch chan *client.SolveStatus) error {
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
					printVertex(vertex, targetConsole)
					logger.Info("Vertex started or cached")
				}
				if vertex.Error != "" {
					if !introducedVertex[vertex.Digest] {
						introducedVertex[vertex.Digest] = true
						printVertex(vertex, targetConsole)
					}
					targetConsole.Printf("ERROR: (%s) %s\n", operation, vertex.Error)
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
				logger = logger.
					With("progress", progress).
					With("name", vs.Name)
				if !introducedVertex[vertex.Digest] {
					introducedVertex[vertex.Digest] = true
					printVertex(vertex, targetConsole)
				}
				logger.Info(vs.ID)
				targetConsole.Printf("%s %d%%\n", vs.ID, progress)
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
				targetConsole.PrintBytes(logLine.Data)
				logger.Info(string(logLine.Data))
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *solver) monitorProgressBasic(ctx context.Context, ch chan *client.SolveStatus) error {
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
				if !introducedVertex[vertex.Digest] {
					if vertex.Cached || vertex.Started != nil {
						introducedVertex[vertex.Digest] = true
						logger.Debug("Vertex")
					}
				}
				if vertex.Error != "" {
					logger.Error(errors.New(vertex.Error))
					if !introducedVertex[vertex.Digest] {
						introducedVertex[vertex.Digest] = true
						printVertex(vertex, targetConsole)
					}
					targetConsole.Printf("ERROR: (%s) %s\n", operation, vertex.Error)
				}
			}
			for _, vs := range ss.Statuses {
				_, found := vertices[vs.Vertex]
				if !found {
					// No logging for internal operations.
					continue
				}
				logger := vertexLoggers[vs.Vertex]
				progress := int32(0)
				if vs.Total != 0 {
					progress = int32(100.0 * float32(vs.Current) / float32(vs.Total))
				}
				if vs.Completed != nil {
					progress = 100
				}
				logger = logger.
					With("progress", int32(progress)).
					With("name", vs.Name)
				logger.Debug(vs.ID)
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

func (s *solver) newSolveOptDocker(img *image.Image, dockerTag string, localDirs map[string]string, w io.WriteCloser) (*client.SolveOpt, error) {
	imgJSON, err := json.Marshal(img)
	if err != nil {
		return nil, errors.Wrap(err, "image json marshal")
	}
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: client.ExporterDocker,
				Attrs: map[string]string{
					"name":                  dockerTag,
					"containerimage.config": string(imgJSON),
				},
				Output: func(_ map[string]string) (io.WriteCloser, error) {
					return w, nil
				},
			},
		},
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
		LocalDirs:           localDirs,
	}, nil
}

func (s *solver) newSolveOptArtifacts(outDir string, localDirs map[string]string) (*client.SolveOpt, error) {
	return &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type:      client.ExporterLocal,
				OutputDir: outDir,
			},
		},
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
		LocalDirs:           localDirs,
	}, nil
}

func (s *solver) newSolveOptSideEffects(localDirs map[string]string) (*client.SolveOpt, error) {
	return &client.SolveOpt{
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
		LocalDirs:           localDirs,
		// CacheImports: []client.CacheOptionsEntry{
		// 	newRegistryCacheOpt("docker.io/earthly/buildkitd:cache"),
		// },
		// CacheExports: []client.CacheOptionsEntry{
		// 	newRegistryCacheOpt("docker.io/earthly/buildkitd:cache"),
		// },
	}, nil
}

func newRegistryCacheOpt(ref string) client.CacheOptionsEntry {
	registryCacheOptAttrs := make(map[string]string)
	registryCacheOptAttrs["ref"] = ref
	registryCacheOptAttrs["mode"] = "max"
	return client.CacheOptionsEntry{
		Type:  "registry",
		Attrs: registryCacheOptAttrs,
	}
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

func shortDigest(d digest.Digest) string {
	return d.Hex()[:12]
}

func loadDockerTar(ctx context.Context, r io.ReadCloser) error {
	// TODO: This is a gross hack - should use proper docker client.
	cmd := exec.CommandContext(ctx, "docker", "load")
	cmd.Stdin = r
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "docker load")
	}
	return nil
}

func pushDockerImage(ctx context.Context, imageName string) error {
	cmd := exec.CommandContext(ctx, "docker", "push", imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrapf(err, "docker push %s", imageName)
	}
	return nil
}
