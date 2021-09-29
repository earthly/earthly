package builder

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/docker/distribution/reference"
	"github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/llbutil"
	"golang.org/x/sync/errgroup"

	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

type manifest struct {
	imageName string
	platform  specs.Platform
}

func platformSpecificImageName(imgName string, platform specs.Platform) (string, error) {
	r, err := reference.ParseNormalizedNamed(imgName)
	if err != nil {
		return "", errors.Wrapf(err, "parse %s", imgName)
	}
	taggedR, ok := reference.TagNameOnly(r).(reference.Tagged)
	if !ok {
		return "", errors.Wrapf(err, "not tagged %s", reference.TagNameOnly(r).String())
	}
	platformTag := llbutil.DockerTagSafe(fmt.Sprintf("%s_%s", taggedR.Tag(), platforms.Format(platform)))
	r2, err := reference.WithTag(r, platformTag)
	if err != nil {
		return "", errors.Wrapf(err, "with tag %s - %s", r.String(), platformTag)
	}
	return reference.FamiliarString(r2), nil
}

func loadDockerManifest(ctx context.Context, console conslogging.ConsoleLogger, parentImageName string, children []manifest) error {
	if len(children) == 0 {
		return errors.Errorf("no images in manifest list for %s", parentImageName)
	}
	// Check if any child has the platform as the default platform (use the first one if none found).
	defaultChild := 0
	for i, child := range children {
		if platforms.Format(child.platform) == platforms.Format(llbutil.DefaultPlatform()) {
			defaultChild = i
			break
		}
	}

	var childImgs []string
	for i, child := range children {
		if i == defaultChild {
			childImgs = append(childImgs, fmt.Sprintf("%s (=%s)", child.imageName, parentImageName))
		} else {
			childImgs = append(childImgs, child.imageName)
		}
	}
	const noteDetail = "Note that when pushing a multi-platform image, " +
		"it is pushed as a single multi-manifest image. " +
		"Separate per-platform image tags are only available locally."
	console.Printf(
		"Image %s is a multi-platform image. The following per-platform images have been produced:\n\t%s\n%s\n",
		parentImageName, strings.Join(childImgs, "\n\t"), noteDetail)

	cmd := exec.CommandContext(ctx, "docker", "tag", children[defaultChild].imageName, parentImageName)
	cmd.Stdout = os.Stderr // Preserve desired output on stdout, all logs to stderr
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "docker tag default platform image")
	}
	return nil
}

func loadDockerTar(ctx context.Context, r io.ReadCloser, console conslogging.ConsoleLogger) error {
	cmd := exec.CommandContext(ctx, "docker", "load")
	cmd.Stdin = r
	output, err := cmd.CombinedOutput()
	if err != nil {
		console.Warnf("%+v output:\n%s\n", cmd.Args, string(output))
		return errors.Wrapf(err, "docker load")
	}
	return nil
}

func dockerPullLocalImages(ctx context.Context, localRegistryAddr string, pullMap map[string]string, console conslogging.ConsoleLogger) error {
	eg, ctx := errgroup.WithContext(ctx)
	for pullName, finalName := range pullMap {
		pn := pullName
		fn := finalName
		eg.Go(func() error {
			return dockerPullLocalImage(ctx, localRegistryAddr, pn, fn, console)
		})
	}
	return eg.Wait()
}

func dockerPullLocalImage(ctx context.Context, localRegistryAddr string, pullName string, finalName string, console conslogging.ConsoleLogger) error {
	fullPullName := fmt.Sprintf("%s/%s", localRegistryAddr, pullName)
	cmd := exec.CommandContext(ctx, "docker", "pull", fullPullName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		console.Warnf("%+v output:\n%s\n", cmd.Args, string(output))
		return errors.Wrapf(err, "docker pull")
	}
	cmd = exec.CommandContext(ctx, "docker", "tag", fullPullName, finalName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		console.Warnf("%+v output:\n%s\n", cmd.Args, string(output))
		return errors.Wrap(err, "docker tag after pull")
	}
	cmd = exec.CommandContext(ctx, "docker", "rmi", fullPullName)
	output, err = cmd.CombinedOutput()
	if err != nil {
		console.Warnf("%+v output:\n%s\n", cmd.Args, string(output))
		return errors.Wrap(err, "docker rmi after pull and retag")
	}
	return nil
}
