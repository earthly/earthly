package containerutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/hashicorp/go-multierror"
	_ "github.com/moby/buildkit/client/connhelper/podmancontainer" // Load "podman-container://" helper.
	"github.com/pkg/errors"
)

type podmanShellFrontend struct {
	*shellFrontend
}

// NewPodmanShellFrontend constructs a new Frontend using the podman binary installed on the host.
// It also ensures that the binary is functional for our needs and collects compatibility information.
func NewPodmanShellFrontend(ctx context.Context, cfg *FrontendConfig) (ContainerFrontend, error) {
	fe := &podmanShellFrontend{
		shellFrontend: &shellFrontend{
			binaryName:              "podman",
			runCompatibilityArgs:    []string{"--security-opt", "unmask=/sys/fs/cgroup"},
			globalCompatibilityArgs: make([]string, 0),
			Console:                 cfg.Console,
		},
	}

	output, err := fe.commandContextOutput(ctx, "info", "--format={{.Host.Security.Rootless}}")
	if err != nil {
		return nil, err
	}

	if output.stderr.Len() > 0 {
		// Only check stdout; since some podman versions less than 3.4 will report warnings about no systemd session,
		// and falling back to cgroupfs. These errors land on stderr. https://github.com/containers/podman/pull/12834

		cfg.Console.VerbosePrintf("Podman logged additional information to stderr:")
		cfg.Console.VerbosePrintf(output.stderr.String())
		cfg.Console.VerbosePrintf("Adding log level compatibility flag for all additional operations.")

		fe.globalCompatibilityArgs = append(fe.globalCompatibilityArgs, "--log-level", "error")
	}

	// Only check stdout here since it may be contaminated with log output detected above.
	trimmedStdOut := strings.TrimSpace(output.stdout.String())
	isRootless, err := strconv.ParseBool(trimmedStdOut)
	if err != nil {
		return nil, errors.Wrapf(err, "info returned invalid value %s", output.string())
	}
	fe.rootless = isRootless

	fe.urls, err = fe.setupAndValidateAddresses(FrontendPodmanShell, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate buildkit URLs")
	}

	return fe, nil
}

func (psf *podmanShellFrontend) Scheme() string {
	return "podman-container"
}

func (psf *podmanShellFrontend) Config() *CurrentFrontend {
	return &CurrentFrontend{
		Setting:      FrontendPodmanShell,
		Binary:       psf.binaryName,
		Type:         FrontendTypeShell,
		FrontendURLs: psf.urls,
	}
}

func (psf *podmanShellFrontend) Information(ctx context.Context) (*FrontendInfo, error) {
	output, err := psf.commandContextOutput(ctx, "info", "--format={{.Host.RemoteSocket.Exists}}")
	if err != nil {
		return nil, err
	}

	hasRemote, err := strconv.ParseBool(output.string())
	if err != nil {
		return nil, errors.Wrapf(err, "info returned invalid value %s", output.string())
	}

	args := []string{"version", "--format=json"}
	if hasRemote {
		args = append([]string{"-r"}, args...)
	}

	output, err = psf.commandContextOutput(ctx, args...)
	if err != nil {
		return nil, err
	}

	type versionInfo struct {
		Version    string
		APIVersion string
		OSArch     string
	}

	type info struct {
		Client versionInfo
		Server versionInfo
	}

	allInfo := info{}
	err = json.Unmarshal([]byte(output.string()), &allInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse version output %s", output.string())
	}

	host := "daemonless"
	if hasRemote {
		output, err = psf.commandContextOutput(ctx, "info", "--format={{.Host.RemoteSocket.Path}}")
		if err != nil {
			return nil, err
		}
		host = string(output.string())
	}

	return &FrontendInfo{
		ClientVersion:    allInfo.Client.Version,
		ClientAPIVersion: allInfo.Client.APIVersion,
		ClientPlatform:   allInfo.Client.OSArch,
		ServerVersion:    allInfo.Server.Version,
		ServerAPIVersion: allInfo.Server.APIVersion,
		ServerPlatform:   allInfo.Server.OSArch,
		ServerAddress:    host,
	}, nil
}

func (psf *podmanShellFrontend) ImagePull(ctx context.Context, refs ...string) error {
	var err error
	for _, ref := range refs {
		args := []string{"pull"}
		if strings.HasPrefix(ref, psf.urls.LocalRegistryHost.Host+"/") {
			// Rather than force users to add an exemption locally in /etc/containers/registries.conf, detect when we are
			// pulling from our own internal registry and manually exempt it from TLS.
			args = append(args, "--tls-verify=false")
		}
		args = append(args, ref)

		_, cmdErr := psf.commandContextOutput(ctx, args...)
		if cmdErr != nil {
			err = multierror.Append(err, cmdErr)
		}
	}

	return err
}

func (psf *podmanShellFrontend) ImageLoadFromFileCommand(filename string) string {
	binary, args := psf.commandContextStrings("pull", fmt.Sprintf("docker-archive:%s", filename))

	all := []string{binary}
	all = append(all, args...)

	return strings.Join(all, " ")
}

func (psf *podmanShellFrontend) ImageLoad(ctx context.Context, images ...io.Reader) error {
	var err error
	for _, image := range images {
		// Write the image to a temp file. This is needed to accommodate some Podman versions between 3.0 and 3.4. Because
		// buildkit creates weird hybrid docker/OCI images, Podman pulls it in as an OCI image and ends up neglecting the
		// in-built image tag. We can get around this by "pulling" a tar file and specifying the format at the CLI. This
		// is more or less what Podman will be doing going forward. For further context, see the linked issues and discussion
		// here: https://github.com/earthly/earthly/issues/1285

		file, tmpErr := os.CreateTemp("", "earthly-podman-load-*")
		if tmpErr != nil {
			err = multierror.Append(err, errors.Wrap(tmpErr, "failed to create temp tarball"))
			continue
		}
		_, copyErr := io.Copy(file, image)
		if copyErr != nil {
			err = multierror.Append(err, errors.Wrapf(tmpErr, "failed to write to %s", file.Name()))
			continue
		}
		defer file.Close()
		defer os.Remove(file.Name())

		output, cmdErr := psf.commandContextOutput(ctx, "pull", fmt.Sprintf("docker-archive:%s", file.Name()))
		if cmdErr != nil {
			err = multierror.Append(err, errors.Wrapf(cmdErr, "image load failed: %s", output.string()))
		}
	}

	return err
}

func (psf *podmanShellFrontend) VolumeInfo(ctx context.Context, volumeNames ...string) (map[string]*VolumeInfo, error) {
	// Older podman versions do no support --format. This means we are stuck parsing the verbose tabular output for compat.
	output, err := psf.commandContextOutput(ctx, "system", "df", "-v")
	if err != nil {
		return nil, err
	}

	idx := strings.Index(output.string(), "Local Volumes space usage:")
	val := output.string()[idx:]
	lines := strings.Split(string(val), "\n")[3:]
	results := map[string]*VolumeInfo{}

	for _, line := range lines {
		lineParts := strings.Fields(line)
		for _, volumeName := range volumeNames {
			// There are three columns. By index:
			// 0 -> name, 1 -> links, 2 -> size
			// There may be straggler lines after due to parsing, ignore them. They will not have enough length.
			// The volume lines are last so we are safe.
			if len(lineParts) == 3 && volumeName == lineParts[0] {
				// Get size
				var bytes uint64
				bytes, parseErr := humanize.ParseBytes(lineParts[2])
				if err != nil {
					err = multierror.Append(err, parseErr)
					break
				}

				// The mountpoint is not included in the df output. Get that from inspect.
				mountpoint, mountpointErr := psf.commandContextOutput(ctx, "volume", "inspect", volumeName, "--format={{.Mountpoint}}")
				if err != nil {
					err = multierror.Append(err, mountpointErr)
					break
				}

				results[volumeName] = &VolumeInfo{
					Name:       volumeName,
					SizeBytes:  bytes,
					Mountpoint: string(mountpoint.string()),
				}
				break
			}
		}
	}

	return results, err
}
