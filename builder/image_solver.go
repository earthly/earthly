package builder

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/earthly/earthly/outmon"
	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/states/image"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/platutil"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
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
			vertexFailureOutput, err = s.sm.MonitorProgress(ctx, ch, "", true)
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

type localRegistryImageSolver struct {
	bkClient        *client.Client
	sm              *outmon.SolverMonitor
	attachables     []session.Attachable
	enttlmnts       []entitlements.Entitlement
	cacheImports    *states.CacheImports
	cacheExport     string
	maxCacheExport  string
	saveInlineCache bool
}

func newLocalRegistryImageSolver(opt Opt, sm *outmon.SolverMonitor) *localRegistryImageSolver {
	return &localRegistryImageSolver{
		sm:              sm,
		bkClient:        opt.BkClient,
		attachables:     opt.Attachables,
		enttlmnts:       opt.Enttlmnts,
		cacheImports:    opt.CacheImports,
		cacheExport:     opt.CacheExport,
		maxCacheExport:  opt.MaxCacheExport,
		saveInlineCache: opt.SaveInlineCache,
	}
}

func (s *localRegistryImageSolver) newSolveOpt(img *image.Image, dockerTag string, onPull pullping.PullCallback) *client.SolveOpt {
	var cacheImports []client.CacheOptionsEntry
	for ci := range s.cacheImports.AsMap() {
		cacheImports = append(cacheImports, newCacheImportOpt(ci))
	}
	var cacheExports []client.CacheOptionsEntry
	if s.cacheExport != "" {
		cacheExports = append(cacheExports, newCacheExportOpt(s.cacheExport, false))
	}
	if s.maxCacheExport != "" {
		cacheExports = append(cacheExports, newCacheExportOpt(s.maxCacheExport, true))
	}
	if s.saveInlineCache {
		cacheExports = append(cacheExports, newInlineCacheOpt())
	}
	return &client.SolveOpt{
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
		CacheExports:        cacheExports,
		Session:             s.attachables,
		AllowedEntitlements: s.enttlmnts,
	}
}

// SolveImage also creates a Docker image but it stores the image using the
// embedded Docker registry in BuildKit. The method returns a string channel to
// which Docker image names written, a close() function that must be called
// after the images have been used, and an error channel to which any errors
// will be sent.
func (s *localRegistryImageSolver) SolveImage(ctx context.Context, mts *states.MultiTarget, dockerTag string) (chan string, func(), chan error, error) {
	var (
		platform  = mts.Final.PlatformResolver.ToLLBPlatform(mts.Final.PlatformResolver.Current())
		saveImage = mts.Final.LastSaveImage()
	)

	var (
		releaseChan = make(chan struct{})
		resultsChan = make(chan string, 1)
		errChan     = make(chan error)
	)

	// This func will be exposed to the caller and must be invoked when the WITH
	// DOCKER command has termintated. It will release any prepared Docker
	// images.
	closer := func() {
		close(releaseChan)
	}

	// This func is execute when the image create/push process is complete.
	onPull := func(ctx context.Context, images []string) error {
		// Send any images created by BuildKit to the caller.
		for _, image := range images {
			resultsChan <- image
		}
		close(resultsChan)
		// Wait for the closer func to be called. This signals that all WITH
		// DOCKER statements have been run and we can release the image
		// resources. When the onPull function returns BK will remove the
		// images.
		<-releaseChan
		return nil
	}

	imgJSON, err := json.Marshal(saveImage.Image)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "image json marshal")
	}

	solveOpt := s.newSolveOpt(saveImage.Image, dockerTag, onPull)
	bf := func(childCtx context.Context, gwClient gwclient.Client) (*gwclient.Result, error) {
		res := gwclient.NewResult()
		ref, err := llbutil.StateToRef(childCtx, gwClient, saveImage.State, true, platutil.NewResolver(platform), s.cacheImports.AsMap())
		if err != nil {
			return nil, errors.Wrap(err, "initial state to ref conversion")
		}

		refKey := fmt.Sprintf("image-%s", dockerTag)
		refPrefix := fmt.Sprintf("ref/%s", refKey)
		res.AddRef(refKey, ref)

		localRegPullID := fmt.Sprintf("sess-%s/%s", gwClient.BuildOpts().SessionID, dockerTag)
		res.AddMeta(fmt.Sprintf("%s/export-image-local-registry", refPrefix), []byte(localRegPullID))
		res.AddMeta(fmt.Sprintf("%s/%s", refPrefix, exptypes.ExporterImageConfigKey), imgJSON)
		res.AddMeta(fmt.Sprintf("%s/image.name", refPrefix), []byte(dockerTag))
		res.AddMeta(fmt.Sprintf("%s/image-index", refPrefix), []byte("0"))

		return res, nil
	}

	// TODO: figure out why the below context is being canceled.
	go func() {
		_, err := s.bkClient.Build(context.TODO(), *solveOpt, "", bf, nil)
		if err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	return resultsChan, closer, errChan, nil
}
