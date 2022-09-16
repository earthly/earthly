package containerutil

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

// This is a stub for use when a proper frontend is not available.
type stubFrontend struct {
	*shellFrontend
}

// ErrFrontendNotInitialized is returned when the frontend is not initialized.
var ErrFrontendNotInitialized = errors.New("frontend (e.g. docker/podman) not initialized")

// NewStubFrontend creates a stubbed frontend. Useful in cases where a frontend could not be detected, but we still need a frontend.
// Examples include earthly/earthly, or integration tests. It is currently only used as a fallback when docker or other frontends are missing.
func NewStubFrontend(ctx context.Context, cfg *FrontendConfig) (ContainerFrontend, error) {
	fe := &stubFrontend{
		shellFrontend: &shellFrontend{},
	}
	fe.shellFrontend.FrontendInformation = fe.Information

	var err error
	fe.urls, err = fe.setupAndValidateAddresses(FrontendStub, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to calculate buildkit URLs")
	}

	return fe, nil
}

func (*stubFrontend) Scheme() string {
	return ""
}
func (*stubFrontend) IsAvailable(ctx context.Context) bool {
	return false
}
func (sf *stubFrontend) Config() *CurrentFrontend {
	return &CurrentFrontend{
		Setting:      FrontendStub,
		FrontendURLs: sf.urls,
	}
}
func (*stubFrontend) Information(ctx context.Context) (*FrontendInfo, error) {
	return &FrontendInfo{}, nil
}
func (*stubFrontend) ContainerInfo(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerInfo, error) {
	return nil, ErrFrontendNotInitialized
}
func (*stubFrontend) ContainerRemove(ctx context.Context, force bool, namesOrIDs ...string) error {
	return ErrFrontendNotInitialized
}
func (*stubFrontend) ContainerStop(ctx context.Context, timeoutSec uint, namesOrIDs ...string) error {
	return ErrFrontendNotInitialized
}
func (*stubFrontend) ContainerLogs(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerLogs, error) {
	return nil, ErrFrontendNotInitialized
}
func (*stubFrontend) ContainerRun(ctx context.Context, containers ...ContainerRun) error {
	return ErrFrontendNotInitialized
}
func (*stubFrontend) ImageInfo(ctx context.Context, refs ...string) (map[string]*ImageInfo, error) {
	return nil, ErrFrontendNotInitialized
}
func (*stubFrontend) ImagePull(ctx context.Context, refs ...string) error {
	return ErrFrontendNotInitialized
}
func (*stubFrontend) ImageRemove(ctx context.Context, force bool, refs ...string) error {
	return ErrFrontendNotInitialized
}
func (*stubFrontend) ImageTag(ctx context.Context, tags ...ImageTag) error {
	return ErrFrontendNotInitialized
}
func (*stubFrontend) ImageLoadFromFileCommand(filename string) string {
	return ""
}
func (*stubFrontend) ImageLoad(ctx context.Context, image ...io.Reader) error {
	return ErrFrontendNotInitialized
}
func (*stubFrontend) VolumeInfo(ctx context.Context, volumeNames ...string) (map[string]*VolumeInfo, error) {
	return nil, ErrFrontendNotInitialized
}
