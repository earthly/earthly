package containerutil

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/pkg/errors"
)

// ContainerFrontend is an interface specifying all the container options Earthly needs to do.
type ContainerFrontend interface {
	Scheme() string

	IsAvaliable(ctx context.Context) bool
	Config() *FrontendConfig
	Information(ctx context.Context) (*FrontendInfo, error)

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

	VolumeInfo(ctx context.Context, volumeNames ...string) (map[string]*VolumeInfo, error)
}

// FrontendForSetting returns a frontend given a setting. This includes automatic detection.
func FrontendForSetting(ctx context.Context, feType string) (ContainerFrontend, error) {
	if feType == FrontendAuto {
		return autodetectFrontend(ctx)
	}

	return frontendIfAvaliable(ctx, feType)
}

func autodetectFrontend(ctx context.Context) (ContainerFrontend, error) {
	if fe, err := frontendIfAvaliable(ctx, FrontendDockerShell); err == nil {
		return fe, nil
	}

	if fe, err := frontendIfAvaliable(ctx, FrontendPodmanShell); err == nil {
		return fe, nil
	}

	return nil, errors.New("failed to autodetect a supported frontend")
}

func frontendIfAvaliable(ctx context.Context, feType string) (ContainerFrontend, error) {
	var newFe func(context.Context) (ContainerFrontend, error)
	switch feType {
	case FrontendDockerShell:
		newFe = NewDockerShellFrontend
	case FrontendPodmanShell:
		newFe = NewPodmanShellFrontend
	default:
		return nil, fmt.Errorf("%s is not a supported container frontend", feType)
	}

	fe, err := newFe(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "%s frontend failed to initalize", feType)
	}
	if !fe.IsAvaliable(ctx) {
		return nil, fmt.Errorf("%s frontend not avaliable", feType)
	}

	return fe, nil
}

func getPlatform() string {
	arch := runtime.GOARCH
	if runtime.GOARCH == "arm" {
		arch = "arm/v7"
	}
	return fmt.Sprintf("linux/%s", arch)
}
