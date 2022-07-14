package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/earthly/earthly/outmon"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/util/gatewaycrafter"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/pullping"
	"github.com/moby/buildkit/util/entitlements"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func newTarImageSolver(opt Opt, sm *outmon.SolverMonitor) *tarImageSolver {
	return &tarImageSolver{
		sm:           sm,
		bkClient:     opt.BkClient,
		attachables:  opt.Attachables,
		enttlmnts:    opt.Enttlmnts,
		cacheImports: opt.CacheImports,
	}
}

type tarImageSolver struct {
	bkClient     *client.Client
	sm           *outmon.SolverMonitor
	attachables  []session.Attachable
	enttlmnts    []entitlements.Entitlement
	cacheImports *states.CacheImports
}

func (s *tarImageSolver) newSolveOpt(img *image.Image, dockerTag string, w io.WriteCloser) (*client.SolveOpt, error) {
	imgJSON, err := json.Marshal(img)
	if err != nil {
		return nil, errors.Wrap(err, "image json marshal")
	}
	var cacheImports []client.CacheOptionsEntry
	for ci := range s.cacheImports.AsMap() {
		cacheImports = append(cacheImports, newCacheImportOpt(ci))
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
		CacheImports:        cacheImports,
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}, nil
}

// SolveImage invokes a BK solve operation to create a Docker image. It is then
// saved to a local tar file.
func (s *tarImageSolver) SolveImage(ctx context.Context, mts *states.MultiTarget, dockerTag string, outFile string, printOutput bool) error {
	platform := mts.Final.PlatformResolver.ToLLBPlatform(mts.Final.PlatformResolver.Current())
	saveImage := mts.Final.LastSaveImage()
	dt, err := saveImage.State.Marshal(ctx, llb.Platform(platform))
	if err != nil {
		return errors.Wrap(err, "state marshal")
	}
	pipeR, pipeW := io.Pipe()
	solveOpt, err := s.newSolveOpt(saveImage.Image, dockerTag, pipeW)
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
		return nil
	})
	var vertexFailureOutput string
	eg.Go(func() error {
		var err error
		if printOutput {
			vertexFailureOutput, err = s.sm.MonitorProgress(ctx, ch, "", true, s.bkClient)
			return err
		}
		// Silent case.
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case _, ok := <-ch:
				if !ok {
					return nil
				}
				// Do nothing - just consume the status updates silently.
			}
		}
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
		<-ctx.Done()
		// Close read pipe on cancels, otherwise the whole thing hangs.
		pipeR.Close()
	}()
	err = eg.Wait()
	if err != nil {
		return NewBuildError(err, vertexFailureOutput)
	}
	return nil
}

type multiImageSolver struct {
	bkClient     *client.Client
	sm           *outmon.SolverMonitor
	attachables  []session.Attachable
	enttlmnts    []entitlements.Entitlement
	cacheImports *states.CacheImports
}

func newMultiImageSolver(opt Opt, sm *outmon.SolverMonitor) *multiImageSolver {
	return &multiImageSolver{
		sm:           sm,
		bkClient:     opt.BkClient,
		attachables:  opt.Attachables,
		enttlmnts:    opt.Enttlmnts,
		cacheImports: opt.CacheImports,
	}
}

// SolveImages uses BuildKit to solve multiple images using a single build
// operation. It stores the images using the embedded Docker registry in
// BuildKit. The method returns a string channel to which Docker image names
// written, a release function that must be called after the images have been
// used, and an error channel to which any errors will be sent.
func (m *multiImageSolver) SolveImages(ctx context.Context, imageDefs []*states.ImageDef) (*states.ImageSolverResults, error) {
	var (
		releaseChan = make(chan struct{})
		resultChan  = make(chan string, len(imageDefs))
		errChan     = make(chan error)
		ret         = &states.ImageSolverResults{
			ResultChan: resultChan,
			ErrChan:    errChan,
			ReleaseFunc: func() {
				close(releaseChan)
			},
		}
	)
	if len(imageDefs) == 0 {
		// Nothing to solve.
		close(resultChan)
		close(errChan)
		return ret, nil
	}

	// This func is executed when the image create/push process is complete.
	onPull := func(ctx context.Context, images []string) error {
		// Send any images created by BuildKit to the caller.
		for _, image := range images {
			resultChan <- image
		}
		close(resultChan)
		// Wait for the closer func to be called. This signals that all WITH
		// DOCKER statements have been run and we can release the image
		// resources. When the onPull function returns BK will remove the
		// images.
		select {
		case <-releaseChan:
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	}

	buildFn := func(childCtx context.Context, gwClient gwclient.Client) (*gwclient.Result, error) {
		gwCrafter := gatewaycrafter.NewGatewayCrafter()

		for i, imageDef := range imageDefs {
			err := m.addRefToResult(childCtx, gwClient, gwCrafter, imageDef, i)
			if err != nil {
				return nil, err
			}
		}

		return gwCrafter.GetResult(), nil
	}

	var (
		statusChan = make(chan *client.SolveStatus)
		doneChan   = make(chan struct{})
	)

	var cacheImports []client.CacheOptionsEntry
	for ci := range m.cacheImports.AsMap() {
		cacheImports = append(cacheImports, newCacheImportOpt(ci))
	}

	solveOpt := &client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type:  client.ExporterEarthly,
				Attrs: map[string]string{},
				// Not used but required in client validation.
				Output: func(map[string]string) (io.WriteCloser, error) {
					return nil, errors.New("not implemented")
				},
				OutputPullCallback: pullping.PullCallback(onPull),
			},
		},
		CacheImports:        cacheImports,
		Session:             m.attachables,
		AllowedEntitlements: m.enttlmnts,
	}

	var vertexFailureOutput string

	go func() {
		_, err := m.bkClient.Build(ctx, *solveOpt, "", buildFn, statusChan)
		if err != nil {
			errChan <- NewBuildError(err, vertexFailureOutput)
		}
		doneChan <- struct{}{}
	}()

	go func() {
		vertexFailureOutput, err := m.sm.MonitorProgress(ctx, statusChan, "", true, m.bkClient)
		if err != nil {
			errChan <- NewBuildError(err, vertexFailureOutput)
		}
		doneChan <- struct{}{}
	}()

	go func() {
		<-doneChan
		<-doneChan
		close(errChan)
	}()

	return ret, nil
}

func (m *multiImageSolver) addRefToResult(ctx context.Context, gwClient gwclient.Client, gwCrafter *gatewaycrafter.GatewayCrafter, imageDef *states.ImageDef, imageIndex int) error {
	saveImage := imageDef.MTS.Final.LastSaveImage()

	if !strings.Contains(imageDef.ImageName, ":") {
		imageDef.ImageName += ":latest"
	}

	ref, err := llbutil.StateToRef(ctx, gwClient, saveImage.State, true, imageDef.MTS.Final.PlatformResolver, m.cacheImports.AsMap())
	if err != nil {
		return errors.Wrap(err, "initial state to ref conversion")
	}

	refPrefix, err := gwCrafter.AddPushImageEntry(ref, imageIndex, imageDef.ImageName, false, false, saveImage.Image, nil)
	if err != nil {
		return err
	}

	localRegPullID := fmt.Sprintf("sess-%s/%s", gwClient.BuildOpts().SessionID, imageDef.ImageName)
	gwCrafter.AddMeta(fmt.Sprintf("%s/export-image-local-registry", refPrefix), []byte(localRegPullID))

	return nil
}
