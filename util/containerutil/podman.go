package containerutil

import (
	"context"
	"encoding/json"
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

func NewPodmanShellFrontend(ctx context.Context) (ContainerFrontend, error) {
	fe := &podmanShellFrontend{
		shellFrontend: &shellFrontend{
			binaryName: "podman",
		},
	}

	output, err := fe.commandContextOutput(ctx, "info", "--format={{.Host.Security.Rootless}}")
	if err != nil {
		return nil, err
	}

	isRootless, err := strconv.ParseBool(output.string())
	if err != nil {
		return nil, errors.Wrapf(err, "info returned invalid value %s", output.string())
	}
	fe.rootless = isRootless

	return fe, nil
}

func (psf *podmanShellFrontend) Scheme() string {
	return "podman-container"
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
	json.Unmarshal([]byte(output.string()), &allInfo)

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
					multierror.Append(err, parseErr)
					break
				}

				// The mountpoint is not included in the df output. Get that from inspect.
				mountpoint, mountpointErr := psf.commandContextOutput(ctx, "volume", "inspect", volumeName, "--format={{.Mountpoint}}")
				if err != nil {
					multierror.Append(err, mountpointErr)
					break
				}

				results[volumeName] = &VolumeInfo{
					Name:       volumeName,
					Size:       bytes,
					Mountpoint: string(mountpoint.string()),
				}
				break
			}
		}
	}

	return results, err
}
