package dockerutil

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/platutil"
	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
)

// Manifest contains docker manifest data
type Manifest struct {
	ImageName string
	Platform  platutil.Platform
}

// LoadDockerManifest loads docker manifests
func LoadDockerManifest(ctx context.Context, console conslogging.ConsoleLogger, fe containerutil.ContainerFrontend, parentImageName string, children []Manifest, platr *platutil.Resolver) error {
	if len(children) == 0 {
		return errors.Errorf("no images in manifest list for %s", parentImageName)
	}
	// Check if any child has the platform as the default platform
	defaultChild := 0
	foundPlatform := false
	for i, child := range children {
		if platr.PlatformEquals(child.Platform, platutil.DefaultPlatform) {
			defaultChild = i
			foundPlatform = true
			break
		}
	}
	if !foundPlatform {
		// fall back to using first defined platform (and display a warning)
		console.Warnf(
			"Failed to find default platform (%s) of multi-platform image %s; defaulting to the first platform type: %s\n",
			platr.Materialize(platutil.DefaultPlatform).String(), parentImageName, children[defaultChild].Platform)
	}

	var childImgs []string
	for i, child := range children {
		if i == defaultChild {
			childImgs = append(childImgs, fmt.Sprintf("%s (=%s)", child.ImageName, parentImageName))
		} else {
			childImgs = append(childImgs, child.ImageName)
		}
	}
	const noteDetail = "Note that when pushing a multi-platform image, " +
		"it is pushed as a single multi-manifest image. " +
		"Separate per-platform image tags are only available locally."
	console.Printf(
		"Image %s is a multi-platform image. The following per-platform images have been produced:\n\t%s\n%s\n",
		parentImageName, strings.Join(childImgs, "\n\t"), noteDetail)

	err := fe.ImageTag(ctx, containerutil.ImageTag{
		SourceRef: children[defaultChild].ImageName,
		TargetRef: parentImageName,
	})
	if err != nil {
		return errors.Wrap(err, "docker tag default platform image")
	}
	return nil
}

// LoadDockerTar loads a docker image via a tar
func LoadDockerTar(ctx context.Context, fe containerutil.ContainerFrontend, r io.ReadCloser) error {
	err := fe.ImageLoad(ctx, r)
	if err != nil {
		return errors.Wrapf(err, "load tar")
	}
	return nil
}

// DockerPullLocalImages pulls a docker image from a local registry
func DockerPullLocalImages(ctx context.Context, fe containerutil.ContainerFrontend, localRegistryAddr string, pullMap map[string]string) error {
	eg, ctx := errgroup.WithContext(ctx)
	for pullName, finalName := range pullMap {
		pn := pullName
		fn := finalName
		eg.Go(func() error {
			return dockerPullLocalImage(ctx, fe, localRegistryAddr, pn, fn)
		})
	}
	return eg.Wait()
}

func dockerPullLocalImage(ctx context.Context, fe containerutil.ContainerFrontend, localRegistryAddr string, pullName string, finalName string) error {
	localRegistryAddr = "localhost:8888"
	fullPullName := fmt.Sprintf("%s/%s", localRegistryAddr, pullName)
	fmt.Printf("PULLING FROM LOCAL REG: %s\n", fullPullName)
	err := fe.ImagePull(ctx, fullPullName)
	fmt.Printf("PULL COMPLETE: %v\n", err)
	if err != nil {
		return errors.Wrapf(err, "image pull")
	}
	err = fe.ImageTag(ctx, containerutil.ImageTag{
		SourceRef: fullPullName,
		TargetRef: finalName,
	})
	if err != nil {
		return errors.Wrap(err, "image tag after pull")
	}
	force := true // Sometimes Docker GCs images automatically (force prevents an error).
	err = fe.ImageRemove(ctx, force, fullPullName)
	if err != nil {
		return errors.Wrap(err, "image rmi after pull and retag")
	}
	return nil
}
