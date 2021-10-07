package containerutil

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
func NewDockerShellFrontend(ctx context.Context) (ContainerFrontend, error) {
	fe := &dockerShellFrontend{
		shellFrontend: &shellFrontend{
			binaryName: "docker",
		},
	}

	output, err := fe.commandContextOutput(ctx, "info", "--format={{.SecurityOptions}}")
	if err != nil {
		return nil, err
	}
	fe.rootless = strings.Contains(output.string(), "rootless")
	fe.userNamespaced = strings.Contains(output.string(), "name=userns")

	if fe.userNamespaced {
		fe.compatibilityArgs = []string{"--userns", "host"}
	}

	return fe, nil
}

func (dsf *dockerShellFrontend) Scheme() string {
	return "docker-container"
}

func (dsf *dockerShellFrontend) Config() *FrontendConfig {
	return &FrontendConfig{
		Setting: FrontendDockerShell,
		Binary:  dsf.binaryName,
		Type:    FrontendTypeShell,
	}
}

func (dsf *dockerShellFrontend) Information(ctx context.Context) (*FrontendInfo, error) {
	output, err := dsf.commandContextOutput(ctx, "version", "--format={{json .}}")
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
	json.Unmarshal([]byte(output.string()), &allInfo)

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
						Size:       bytes,
						Mountpoint: volumeInfo.Mountpoint,
					}
				}
				break
			}
		}
	}

	return results, err
}
