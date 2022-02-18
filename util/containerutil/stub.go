package containerutil

import (
	"context"
	"io"

	"github.com/pkg/errors"
)

// This is a stub for use in internal testing when its too much effort to provide a legitimate backend.
// Should never be used IRL.
type stubFrontend struct {
	*shellFrontend
}

// NewStubFrontend creates a stubbed frontend. Useful in cases where a frontend could not be detected, but we still need a frontend.
// Examples include earthly/earthly, or integration tests. It is currently only used as a fallback when docker or other frontends are missing.
func NewStubFrontend(ctx context.Context, cfg *FrontendConfig) (ContainerFrontend, error) {
	fe := &stubFrontend{
		shellFrontend: &shellFrontend{},
	}

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
func (*stubFrontend) IsAvaliable(ctx context.Context) bool {
	return true
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
	return make(map[string]*ContainerInfo), nil
}
func (*stubFrontend) ContainerRemove(ctx context.Context, force bool, namesOrIDs ...string) error {
	return nil
}
func (*stubFrontend) ContainerStop(ctx context.Context, timeoutSec uint, namesOrIDs ...string) error {
	return nil
}
func (*stubFrontend) ContainerLogs(ctx context.Context, namesOrIDs ...string) (map[string]*ContainerLogs, error) {
	return make(map[string]*ContainerLogs), nil
}
func (*stubFrontend) ContainerRun(ctx context.Context, containers ...ContainerRun) error {
	return nil
}
func (*stubFrontend) ImageInfo(ctx context.Context, refs ...string) (map[string]*ImageInfo, error) {
	return make(map[string]*ImageInfo), nil
}
func (*stubFrontend) ImagePull(ctx context.Context, refs ...string) error {
	return nil
}
func (*stubFrontend) ImageRemove(ctx context.Context, force bool, refs ...string) error {
	return nil
}
func (*stubFrontend) ImageTag(ctx context.Context, tags ...ImageTag) error {
	return nil
}
func (*stubFrontend) ImageLoad(ctx context.Context, image ...io.Reader) error {
	return nil
}
func (*stubFrontend) VolumeInfo(ctx context.Context, volumeNames ...string) (map[string]*VolumeInfo, error) {
	return make(map[string]*VolumeInfo), nil
}
