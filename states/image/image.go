package image

import (
	"github.com/earthly/earthly/util/llbutil"
	"github.com/moby/buildkit/exporter/containerimage/image"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// Image is a partial of the standard Image struct defined as part of the image opencontainers spec
// at https://github.com/opencontainers/image-spec/blob/master/specs-go/v1/config.go#L82
type Image struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
	Config       Config `json:"config"`
}

// NewImage returns a new image.
func NewImage() *Image {
	return &Image{
		Config: Config{
			ImageConfig: specs.ImageConfig{
				ExposedPorts: make(map[string]struct{}),
				Labels:       make(map[string]string),
				Volumes:      make(map[string]struct{}),
				Env:          []string{"PATH=" + llbutil.DefaultPathEnv},
				WorkingDir:   "/",
			},
		},
	}
}

// Clone creates a copy of the image.
func (img *Image) Clone() *Image {
	if img == nil {
		return NewImage()
	}
	clone := &Image{
		Architecture: img.Architecture,
		OS:           img.OS,
		Config: Config{
			ImageConfig: specs.ImageConfig{
				User:         img.Config.User,
				Env:          make([]string, len(img.Config.Env)),
				Entrypoint:   make([]string, len(img.Config.Entrypoint)),
				Cmd:          make([]string, len(img.Config.Cmd)),
				WorkingDir:   img.Config.WorkingDir,
				StopSignal:   img.Config.StopSignal,
				ExposedPorts: make(map[string]struct{}),
				Volumes:      make(map[string]struct{}),
				Labels:       make(map[string]string),
			},
		},
	}
	if img.Config.Healthcheck != nil {
		clone.Config.Healthcheck = &image.HealthConfig{
			Test:        make([]string, len(img.Config.Healthcheck.Test)),
			Interval:    img.Config.Healthcheck.Interval,
			Timeout:     img.Config.Healthcheck.Timeout,
			StartPeriod: img.Config.Healthcheck.StartPeriod,
			Retries:     img.Config.Healthcheck.Retries,
		}
		copy(clone.Config.Healthcheck.Test, img.Config.Healthcheck.Test)
	}
	copy(clone.Config.Env, img.Config.Env)
	copy(clone.Config.Entrypoint, img.Config.Entrypoint)
	copy(clone.Config.Cmd, img.Config.Cmd)
	for k, v := range img.Config.ExposedPorts {
		clone.Config.ExposedPorts[k] = v
	}
	for k, v := range img.Config.Volumes {
		clone.Config.Volumes[k] = v
	}
	for k, v := range img.Config.Labels {
		clone.Config.Labels[k] = v
	}
	return clone
}

// Config is a docker compatible config for an image.
type Config struct {
	specs.ImageConfig

	Healthcheck *image.HealthConfig `json:",omitempty"`
}
