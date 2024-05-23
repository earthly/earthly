package containerutil

import (
	"context"
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/earthly/earthly/conslogging"
)

// ContainerFrontend is an interface specifying all the container options Earthly needs to do.
type ContainerFrontend interface {
	Scheme() string

	IsAvailable(ctx context.Context) bool
	Config() *CurrentFrontend
	Information(ctx context.Context) (*FrontendInfo, error)

	ContainerList(ctx context.Context) ([]*ContainerInfo, error)
	ContainerInfo(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerInfo, error)
	ContainerRemove(ctx context.Context, force bool, namesOrIDs ...string) error
	ContainerStop(ctx context.Context, timeoutSec uint, namesOrIDs ...string) error
	ContainerLogs(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerLogs, error)
	ContainerRun(ctx context.Context, containers ...ContainerRun) error

	ImageInfo(ctx context.Context, refs ...string) (map[string]*ImageInfo, error)
	ImagePull(ctx context.Context, refs ...string) error
	ImageRemove(ctx context.Context, force bool, refs ...string) error
	ImageTag(ctx context.Context, tags ...ImageTag) error
	ImageLoad(ctx context.Context, image ...io.Reader) error
	ImageLoadFromFileCommand(filename string) string

	VolumeInfo(ctx context.Context, volumeNames ...string) (map[string]*VolumeInfo, error)
}

// FrontendConfig is the configuration needed to bring up a given frontend. Includes logging and needed information to
// calculate URLs to reach the container.
type FrontendConfig struct {
	BuildkitHostCLIValue  string
	BuildkitHostFileValue string

	LocalRegistryHostFileValue string

	LocalContainerName string
	DefaultPort        int

	Console conslogging.ConsoleLogger
}

// FrontendForSetting returns a frontend given a setting. This includes automatic detection.
func FrontendForSetting(ctx context.Context, feType string, cfg *FrontendConfig) (ContainerFrontend, error) {
	if feType == FrontendAuto {
		return autodetectFrontend(ctx, cfg)
	}

	return frontendIfAvailable(ctx, feType, cfg)
}

func autodetectFrontend(ctx context.Context, cfg *FrontendConfig) (ContainerFrontend, error) {
	var errs error

	for _, feType := range []string{
		FrontendDockerShell,
		FrontendPodmanShell,
	} {
		fe, err := frontendIfAvailable(ctx, feType, cfg)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		if dsf, ok := fe.(*dockerShellFrontend); ok && dsf.likelyPodman {
			// Docker CLI works, but it's likely podman making itself available via docker CLI.
			continue
		}
		return fe, nil
	}
	return nil, errors.Wrapf(errs, "failed to autodetect a supported frontend")
}

func frontendIfAvailable(ctx context.Context, feType string, cfg *FrontendConfig) (ContainerFrontend, error) {
	var newFe func(context.Context, *FrontendConfig) (ContainerFrontend, error)
	switch feType {
	case FrontendDockerShell:
		newFe = NewDockerShellFrontend
	case FrontendPodmanShell:
		newFe = NewPodmanShellFrontend
	default:
		return nil, fmt.Errorf("%s is not a supported container frontend", feType)
	}

	fe, err := newFe(ctx, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "%s frontend failed to initialize", feType)
	}
	if !fe.IsAvailable(ctx) {
		return nil, fmt.Errorf("%s frontend not available", feType)
	}

	return fe, nil
}
