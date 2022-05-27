package earthfile2llb

import (
	"context"
	"fmt"
	"sync"

	"github.com/earthly/earthly/states"
	"github.com/earthly/earthly/util/gatewaycrafter"
	"github.com/earthly/earthly/util/llbutil"
	"github.com/earthly/earthly/util/llbutil/pllb"
	"github.com/earthly/earthly/util/syncutil/semutil"
	"github.com/earthly/earthly/util/syncutil/serrgroup"

	gwclient "github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

type saveImageWaitItem struct {
	c    *Converter
	si   states.SaveImage
	push bool
}

type stateWaitItem struct {
	c     *Converter
	state *pllb.State
}

// waitItem should be either saveImageWaitItem or stateWaitItem
type waitItem interface {
}

type waitBlock struct {
	items []waitItem
	mu    sync.Mutex
}

func newWaitBlock() *waitBlock {
	return &waitBlock{}
}

func (wb *waitBlock) addSaveImage(si states.SaveImage, c *Converter, push bool) {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	item := saveImageWaitItem{
		c:    c,
		si:   si,
		push: push,
	}
	wb.items = append(wb.items, &item)
}

func (wb *waitBlock) addState(state *pllb.State, c *Converter) {
	wb.mu.Lock()
	defer wb.mu.Unlock()
	item := stateWaitItem{
		c:     c,
		state: state,
	}
	wb.items = append(wb.items, &item)
}

func (wb *waitBlock) wait(ctx context.Context) error {
	wb.mu.Lock()
	defer wb.mu.Unlock()

	errGroup, ctx := serrgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return wb.saveImages(ctx)
	})
	errGroup.Go(func() error {
		return wb.waitStates(ctx)
	})
	return errGroup.Wait()
}

func (wb *waitBlock) saveImages(ctx context.Context) error {
	isMultiPlatform := make(map[string]bool)        // DockerTag -> bool
	noManifestListImgs := make(map[string]struct{}) // set based on DockerTag
	platformImgNames := make(map[string]bool)

	imageWaitItems := []*saveImageWaitItem{}
	for _, item := range wb.items {
		saveImage, ok := item.(*saveImageWaitItem)
		if !ok {
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
				isMultiPlatform[saveImage.si.DockerTag] = true // do I need to count for previsouly seen?
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

	refID := 0
	for _, item := range imageWaitItems {
		if !item.push {
			continue
		}

		ref, err := llbutil.StateToRef(
			ctx, item.c.opt.GwClient, item.si.State, item.c.opt.NoCache,
			item.c.platr, item.c.opt.CacheImports.AsMap())
		if err != nil {
			return errors.Wrapf(err, "failed to solve image required for %s", item.si.DockerTag)
		}

		var platformBytes []byte
		if isMultiPlatform[item.si.DockerTag] {
			platformBytes = []byte(item.si.Platform.String())
			platformImgName, err := llbutil.PlatformSpecificImageName(item.si.DockerTag, item.si.Platform)
			if err != nil {
				return err
			}

			if item.si.CheckDuplicate && item.si.DockerTag != "" {
				if _, found := platformImgNames[platformImgName]; found {
					return errors.Errorf(
						"image %s is defined multiple times for the same platform (%s)",
						item.si.DockerTag, platformImgName)
				}
				platformImgNames[platformImgName] = true
			}
		}

		_, err = gwCrafter.AddPushImageEntry(ref, refID, item.si.DockerTag, true, item.si.InsecurePush, item.si.Image, platformBytes)
		if err != nil {
			return err
		}
		refID++
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
