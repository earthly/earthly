package earthfile2llb

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/domain"
	"github.com/earthly/earthly/util/dockerutil"
	"github.com/earthly/earthly/util/gatewaycrafter"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/saveartifactlocally"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/syncutil/serrgroup"
	"github.com/earthly/earthly/util/waitutil"

	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

type waitBlock struct {
	items     []waitutil.WaitItem
	seenItems map[waitutil.WaitItem]struct{}
	mu        sync.Mutex

	// used for short-circuiting
	called            bool
	pushCalled        bool
	localExportCalled bool
}

func newWaitBlock() *waitBlock {
	return &waitBlock{
		seenItems: map[waitutil.WaitItem]struct{}{},
	}
}

func (wb *waitBlock) SetDoSaves() {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	for _, wi := range wb.items {
		wi.SetDoSave()
	}
}

func (wb *waitBlock) SetDoPushes() {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	for _, wi := range wb.items {
		wi.SetDoPush()
	}
}

func (wb *waitBlock) AddItem(item waitutil.WaitItem) {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	_, exists := wb.seenItems[item]
	if exists {
		return
	}
	wb.seenItems[item] = struct{}{}
	wb.items = append(wb.items, item)
}

func (wb *waitBlock) Wait(ctx context.Context, push, localExport bool) error {
	wb.mu.Lock()
	defer wb.mu.Unlock()

	shortCircuit := wb.called
	wb.called = true
	if push && !wb.pushCalled {
		shortCircuit = false
		wb.pushCalled = true
	}
	if localExport && !wb.localExportCalled {
		shortCircuit = false
		wb.localExportCalled = true
	}
	if shortCircuit {
		return nil
	}

	errGroup, ctx := serrgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return wb.saveImages(ctx, push, localExport)
	})
	if localExport {
		errGroup.Go(func() error {
			return wb.saveArtifactLocal(ctx)
		})
	}
	errGroup.Go(func() error {
		return wb.waitStates(ctx)
	})
	return errGroup.Wait()
}

func (wb *waitBlock) saveImages(ctx context.Context, pushesAllowed, localExportsAllowed bool) error {
	isMultiPlatform := make(map[string]bool)        // DockerTag -> bool
	noManifestListImgs := make(map[string]struct{}) // set based on DockerTag
	platformImgNames := make(map[string]bool)
	singPlatImgNames := make(map[string]bool) // ensure that these are unique

	imageWaitItems := []*saveImageWaitItem{}
	for _, item := range wb.items {
		saveImage, ok := item.(*saveImageWaitItem)
		if !ok {
			continue
		}

		if !saveImage.doPush && !saveImage.localExport {
			continue
		}

		if hasPlatform, ok := isMultiPlatform[saveImage.si.DockerTag]; ok {
			if saveImage.si.HasPlatform != hasPlatform {
				return fmt.Errorf("SAVE IMAGE %s is defined multiple times, but not all commands defined a --platform value", saveImage.si.DockerTag)
			}
			if !hasPlatform {
				return fmt.Errorf("SAVE IMAGE %s was already declared (none had --platform values)", saveImage.si.DockerTag)
			}
			if _, found := noManifestListImgs[saveImage.si.DockerTag]; found {
				return fmt.Errorf("cannot save image %s defined multiple times, but declared as SAVE IMAGE --no-manifest-list", saveImage.si.DockerTag)
			}
		}

		if saveImage.si.HasPlatform {
			// SAVE IMAGE was called with a --platform value
			if saveImage.si.NoManifestList {
				noManifestListImgs[saveImage.si.DockerTag] = struct{}{}
				isMultiPlatform[saveImage.si.DockerTag] = false
			} else {
				isMultiPlatform[saveImage.si.DockerTag] = true // do I need to count for previously seen?
			}
		} else {
			isMultiPlatform[saveImage.si.DockerTag] = false
		}
		imageWaitItems = append(imageWaitItems, saveImage)
	}
	if len(imageWaitItems) == 0 {
		return nil
	}

	gwCrafter := gatewaycrafter.NewGatewayCrafter()

	// these are used to pass manifest data to the onImage function in builder.go; this only applies to non-local-registry exports
	var tarImagesInWaitBlockRefPrefixes []string
	var tarImagesInWaitBlock []string

	refID := 0
	for _, item := range imageWaitItems {
		sessionID := item.c.opt.GwClient.BuildOpts().SessionID
		exportCoordinator := item.c.opt.ExportCoordinator
		ref, err := llbutil.StateToRef(
			ctx, item.c.opt.GwClient, item.si.State, item.c.opt.NoCache,
			item.c.platr, item.c.opt.CacheImports.AsSlice())
		if err != nil {
			return errors.Wrapf(err, "failed to solve image required for %s", item.si.DockerTag)
		}

		var platformBytes []byte
		var platformImgName string
		if isMultiPlatform[item.si.DockerTag] {
			platformBytes = []byte(item.si.Platform.String())
			platformImgName, err = llbutil.PlatformSpecificImageName(item.si.DockerTag, item.si.Platform)
			if err != nil {
				return err
			}

			if item.si.CheckDuplicate && item.si.DockerTag != "" {
				if _, found := platformImgNames[platformImgName]; found {
					return errors.Errorf(
						"image %s is defined multiple times for the same platform (%s)",
						item.si.DockerTag, item.si.Platform.String())
				}
				platformImgNames[platformImgName] = true
			}
		} else {
			if item.si.CheckDuplicate && item.si.DockerTag != "" {
				if _, found := singPlatImgNames[item.si.DockerTag]; found {
					return errors.Errorf(
						"image %s is defined multiple times for the same default platform",
						item.si.DockerTag)
				}
				singPlatImgNames[item.si.DockerTag] = true
			}
		}

		refPrefix, err := gwCrafter.AddPushImageEntry(ref, refID, item.si.DockerTag, item.doPush, item.si.InsecurePush, item.si.Image, platformBytes)
		if err != nil {
			return err
		}
		refID++

		if item.localExport {
			if isMultiPlatform[item.si.DockerTag] {
				// local docker instance does not support multi-platform images, so we must create a new entry and set it to the platformImgName
				refPrefix, err := gwCrafter.AddPushImageEntry(ref, refID, platformImgName, false, false, item.si.Image, nil)
				if err != nil {
					return err
				}

				exportCoordinatorImageID := exportCoordinator.AddImage(sessionID, item.si.DockerTag, &dockerutil.Manifest{
					ImageName: platformImgName,
					Platform:  item.si.Platform,
				})

				if item.c.opt.UseLocalRegistry {
					gwCrafter.AddMeta(fmt.Sprintf("%s/export-image-local-registry", refPrefix), []byte(exportCoordinatorImageID))
				} else {
					gwCrafter.AddMeta(fmt.Sprintf("%s/export-image", refPrefix), []byte("true"))
					gwCrafter.AddMeta(fmt.Sprintf("%s/export-image-manifest-key", refPrefix), []byte(exportCoordinatorImageID))
					tarImagesInWaitBlockRefPrefixes = append(tarImagesInWaitBlockRefPrefixes, refPrefix)
					tarImagesInWaitBlock = append(tarImagesInWaitBlock, exportCoordinatorImageID)
				}
				refID++
			} else {
				if item.c.opt.UseLocalRegistry {
					exportCoordinatorImageID := exportCoordinator.AddImage(sessionID, item.si.DockerTag, nil)
					gwCrafter.AddMeta(fmt.Sprintf("%s/export-image-local-registry", refPrefix), []byte(exportCoordinatorImageID))
				} else {
					gwCrafter.AddMeta(fmt.Sprintf("%s/export-image", refPrefix), []byte("true"))
				}
			}
			exportCoordinator.AddLocalOutputSummary(item.c.target.String(), item.si.DockerTag, item.c.mts.Final.ID)
		}
	}
	if len(tarImagesInWaitBlockRefPrefixes) != 0 {
		waitFor := strings.Join(tarImagesInWaitBlock, " ")
		// the wait-for entry is used to know when all multiplatform images have been exported, thus making it safe to load manifests
		for _, refPrefix := range tarImagesInWaitBlockRefPrefixes {
			gwCrafter.AddMeta(fmt.Sprintf("%s/export-image-wait-for", refPrefix), []byte(waitFor))
		}
	}

	if len(imageWaitItems) == 0 {
		panic("saveImagesWaitItem should never have been created with zero converters")
	}
	gatewayClient := imageWaitItems[0].c.opt.GwClient // could be any converter's gwClient (they should app be the same)

	refs, metadata := gwCrafter.GetRefsAndMetadata()
	err := gatewayClient.Export(ctx, gwclient.ExportRequest{
		Refs:     refs,
		Metadata: metadata,
	})
	if err != nil {
		return errors.Wrap(err, "failed to SAVE IMAGE")
	}
	return nil
}

func (wb *waitBlock) waitStates(ctx context.Context) error {

	stateItems := []*stateWaitItem{}

	for _, item := range wb.items {
		stateItem, ok := item.(*stateWaitItem)
		if !ok {
			continue
		}
		stateItems = append(stateItems, stateItem)
	}

	if len(stateItems) == 0 {
		return nil
	}

	// all converters have the same semaphore
	sharedParallelism := stateItems[0].c.opt.Parallelism

	// This semaphore ensures that there is at least one thread allowed to progress,
	// even if parallelism is completely starved.
	sem := semutil.NewMultiSem(sharedParallelism, semutil.NewWeighted(1))

	errGroup, ctx := serrgroup.WithContext(ctx)
	for _, item := range stateItems {
		item := item // must create a new instance here for use in the threaded function
		errGroup.Go(func() error {
			rel, err := sem.Acquire(ctx, 1)
			if err != nil {
				return errors.Wrapf(err, "acquiring parallelism semaphore during waitStates for %s", item.c.target.String())
			}
			defer rel()
			return item.c.forceExecution(ctx, *item.state, item.c.platr)
		})
	}
	return errGroup.Wait()
}

type saveArtifactLocalEntry struct {
	artifact    domain.Artifact
	artifactDir string
	destPath    string
	ifExists    bool
	salt        string
}

func (wb *waitBlock) saveArtifactLocal(ctx context.Context) error {
	gwCrafter := gatewaycrafter.NewGatewayCrafter()

	var gatewayClient gwclient.Client
	var console conslogging.ConsoleLogger
	var exportCoordinator *gatewaycrafter.ExportCoordinator
	artifacts := []saveArtifactLocalEntry{}

	for refID, item := range wb.items {
		saveLocalItem, ok := item.(*saveArtifactLocalWaitItem)
		if !ok {
			continue
		}

		c := saveLocalItem.c
		i := saveLocalItem.saveLocal.Index
		gatewayClient = c.opt.GwClient
		console = c.opt.Console
		exportCoordinator = c.opt.ExportCoordinator

		state := c.mts.Final.SeparateArtifactsState[i]

		ref, err := llbutil.StateToRef(ctx, c.opt.GwClient, state, c.opt.NoCache, c.platr, c.opt.CacheImports.AsSlice())
		if err != nil {
			return err
		}

		artifact := domain.Artifact{
			Target:   c.target,
			Artifact: saveLocalItem.saveLocal.ArtifactPath,
		}

		dirID, err := gwCrafter.AddSaveArtifactLocal(ref, refID, artifact.String(), saveLocalItem.saveLocal.ArtifactPath, saveLocalItem.saveLocal.DestPath)
		if err != nil {
			return err
		}
		c.opt.LocalArtifactWhiteList.Add(saveLocalItem.saveLocal.DestPath)

		outDir, err := c.opt.TempEarthlyOutDir()
		if err != nil {
			return err
		}
		artifacts = append(artifacts, saveArtifactLocalEntry{
			artifact:    artifact,
			artifactDir: filepath.Join(outDir, fmt.Sprintf("index-%s", dirID)),
			destPath:    saveLocalItem.saveLocal.DestPath,
			ifExists:    saveLocalItem.saveLocal.IfExists,
			salt:        c.mts.Final.ID,
		})

	}

	refs, metadata := gwCrafter.GetRefsAndMetadata()
	if len(refs) == 0 {
		if len(metadata) != 0 {
			panic("metadata should always be empty when refs is empty")
		}
		return nil
	}

	err := gatewayClient.Export(ctx, gwclient.ExportRequest{
		Refs:     refs,
		Metadata: metadata,
	})
	if err != nil {
		return err
	}

	for _, entry := range artifacts {
		err = saveartifactlocally.SaveArtifactLocally(
			ctx, exportCoordinator, console, entry.artifact, entry.artifactDir, entry.destPath, entry.salt, entry.ifExists)
		if err != nil {
			return err
		}

	}
	return nil
}
