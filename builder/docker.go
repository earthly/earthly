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
	"github.com/earthly/earthly/llbutil"

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
	console = console.WithPrefix(parentImageName)
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
		"%s is a multi-platform image. The following per-platform images have been produced:\n\t%s\n%s\n",
		parentImageName, strings.Join(childImgs, "\n\t"), noteDetail)

	cmd := exec.CommandContext(ctx, "docker", "tag", children[defaultChild].imageName, parentImageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "docker tag default platform image")
	}
	return nil
}

// TODO: This doesn't work with vanilla docker installations. Not currently used.
func loadDockerManifestExperimental(ctx context.Context, parentImageName string, children []manifest) error {
	createArgs := []string{"manifest", "create", parentImageName}
	for _, child := range children {
		createArgs = append(createArgs, child.imageName)
	}
	cmd := exec.CommandContext(ctx, "docker", createArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "docker manifest create")
	}

	for _, child := range children {
		annotateArgs := []string{
			"manifest", "annotate",
			"--os", child.platform.OS,
			"--arch", child.platform.Architecture,
			"--variant", child.platform.Variant,
			"--os-version", child.platform.OSVersion,
			"--os-features", strings.Join(child.platform.OSFeatures, ","),
		}
		cmd := exec.CommandContext(ctx, "docker", annotateArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return errors.Wrap(err, "docker manifest annotate")
		}
	}
	return nil
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
