package containerutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alessio/shellescape"
	"github.com/dustin/go-humanize"
	"github.com/hashicorp/go-multierror"
	_ "github.com/moby/buildkit/client/connhelper/dockercontainer" // Load "docker-container://" helper.
	"github.com/pkg/errors"
)

type dockerShellFrontend struct {
	*shellFrontend
	userNamespaced bool
}

// NewDockerShellFrontend constructs a new Frontend using the docker binary installed on the host.
// It also ensures that the binary is functional for our needs and collects compatibility information.
func NewDockerShellFrontend(ctx context.Context, cfg *FrontendConfig) (ContainerFrontend, error) {
	fe := &dockerShellFrontend{
		shellFrontend: &shellFrontend{
			binaryName:              "docker",
			runCompatibilityArgs:    make([]string, 0),
			globalCompatibilityArgs: make([]string, 0),
			Console:                 cfg.Console,
		},
	}
	// TODO: Find a cleaner way to pass down information to the shellFrontend
	fe.FrontendInformation = fe.Information

	// running `docker info --format={{.SecurityOptions}}` results in a panic() when docker is not running.
	// To workaround this issue, first we run `docker info` to test docker is running, then again with the
	// `--format` option.
	// This is to prevent displaying panic() errors to our users (even though the panic() occurred in the
	// docker cli binary and not earthly).
	_, err := fe.commandContextOutputWithRetry(ctx, 10, 10*time.Second, "info")
	if err != nil {
		return nil, err
	}

	output, err := fe.commandContextOutputWithRetry(ctx, 10, 10*time.Second, "info", "--format={{.SecurityOptions}}")
	if err != nil {
		return nil, err
	}
	fe.rootless = strings.Contains(output.string(), "rootless")
	fe.userNamespaced = strings.Contains(output.string(), "name=userns")
	if fe.userNamespaced {
		fe.runCompatibilityArgs = []string{"--userns", "host"}
	}
	fe.urls, err = fe.setupAndValidateAddresses(FrontendDockerShell, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate buildkit URLs")
	}

	output, err = fe.commandContextOutputWithRetry(ctx, 10, 10*time.Second, "info", "--format={{.DockerRootDir}}")
	if err != nil {
		return nil, err
	}
	if strings.Contains(output.string(), "/var/lib/containers/storage") {
		// Likely podman making itself available via the docker CLI.
		fe.shellFrontend.likelyPodman = true
	}

	return fe, nil
}

func (dsf *dockerShellFrontend) Scheme() string {
	return "docker-container"
}

func (dsf *dockerShellFrontend) Config() *CurrentFrontend {
	return &CurrentFrontend{
		Setting:      FrontendDockerShell,
		Binary:       dsf.binaryName,
		Type:         FrontendTypeShell,
		FrontendURLs: dsf.urls,
	}
}

func (dsf *dockerShellFrontend) Information(ctx context.Context) (*FrontendInfo, error) {
	output, err := dsf.commandContextOutputWithRetry(ctx, 10, 10*time.Second, "version", "--format={{json .}}")
	if err != nil {
		return nil, err
	}

	type versionInfo struct {
		Version    string
		APIVersion string
		OS         string
		Arch       string
	}

	type info struct {
		Client versionInfo
		Server versionInfo
	}

	allInfo := info{}
	err = json.Unmarshal([]byte(output.string()), &allInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse docker version output")
	}

	host, exists := os.LookupEnv("DOCKER_HOST")
	if !exists {
		host = "/var/run/docker.sock"
	}

	return &FrontendInfo{
		ClientVersion:    allInfo.Client.Version,
		ClientAPIVersion: allInfo.Client.APIVersion,
		ClientPlatform:   fmt.Sprintf("%s/%s", allInfo.Client.OS, allInfo.Client.Arch),
		ServerVersion:    allInfo.Server.Version,
		ServerAPIVersion: allInfo.Server.APIVersion,
		ServerPlatform:   fmt.Sprintf("%s/%s", allInfo.Server.OS, allInfo.Server.Arch),
		ServerAddress:    host,
	}, nil
}

func (dsf *dockerShellFrontend) ContainerInfo(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerInfo, error) {
	results, err := dsf.shellFrontend.ContainerInfo(ctx, namesOrIDs...)
	if err != nil {
		return nil, err
	}

	for _, v := range results {
		// Docker prepends a `\`. This is as intended, according to docker; but unexpected in our
		// case. So remove it. If the status is missing, it was passed through so do not remove.
		if v.Status != StatusMissing {
			v.Name = v.Name[1:]
		}
	}

	return results, nil
}

func (dsf *dockerShellFrontend) ImagePull(ctx context.Context, refs ...string) error {
	var err error
	for _, ref := range refs {
		_, cmdErr := dsf.commandContextOutput(ctx, "pull", ref)
		if cmdErr != nil {
			err = multierror.Append(err, cmdErr)
		}
	}

	return err
}

func (dsf *dockerShellFrontend) ImageLoadFromFileCommand(filename string) string {
	binary, args := dsf.commandContextStrings("load")

	all := []string{binary}
	all = append(all, args...)

	return fmt.Sprintf("cat %s | %s", shellescape.Quote(filename), strings.Join(all, " "))
}

func (dsf *dockerShellFrontend) ImageLoad(ctx context.Context, images ...io.Reader) error {
	var err error
	args := append(dsf.globalCompatibilityArgs, "load")
	for _, image := range images {
		// Do not use the wrapper to allow the image to come in on stdin
		cmd := exec.CommandContext(ctx, dsf.binaryName, args...)
		cmd.Stdin = image
		output, cmdErr := cmd.CombinedOutput()
		if cmdErr != nil {
			err = multierror.Append(err, errors.Wrapf(cmdErr, "image load failed: %s", string(output)))
		}
	}

	return err
}

func (dsf *dockerShellFrontend) VolumeInfo(ctx context.Context, volumeNames ...string) (map[string]*VolumeInfo, error) {
	// Ignore the error. This is because one or more of the provided names could be missing.
	// This allows for Info to report that the volume itself is missing.
	output, _ := dsf.commandContextOutput(ctx, "system", "df", "-v", "--format={{json  .}}")

	results := map[string]*VolumeInfo{}
	for _, name := range volumeNames {
		// Preinitialize all as missing. It will get overwritten when we encounter a real one from the actual output.
		results[name] = &VolumeInfo{Name: name}
	}

	// Anonymous struct to just pick out what we need
	volumeInfos := struct {
		Volumes []struct {
			Name       string `json:"Name"`
			Size       string `json:"Size"`
			Mountpoint string `json:"Mountpoint"`
		} `json:"Volumes"`
	}{}
	err := json.Unmarshal([]byte(output.stdout.String()), &volumeInfos)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode docker volume info for %v", volumeNames)
	}

	for _, name := range volumeNames {
		for _, volumeInfo := range volumeInfos.Volumes {
			if name == volumeInfo.Name {
				bytes, parseErr := humanize.ParseBytes(volumeInfo.Size)
				if parseErr != nil {
					err = multierror.Append(err, parseErr)
				} else {
					results[name] = &VolumeInfo{
						Name:       volumeInfo.Name,
						SizeBytes:  bytes,
						Mountpoint: volumeInfo.Mountpoint,
					}
				}
				break
			}
		}
	}

	return results, err
}
